package interviewHttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	"joblessness/haha/models/base"
	"joblessness/haha/summary/interfaces"
	"joblessness/haha/utils/chat"
	"net/http"
	"strconv"
	"time"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type Handler struct {
	useCase  interviewInterfaces.InterviewUseCase
	upGrader *websocket.Upgrader
	room     chat.Room
}

func NewHandler(useCase interviewInterfaces.InterviewUseCase, room chat.Room) *Handler {
	handler := &Handler{
		useCase: useCase,
		upGrader: &websocket.Upgrader{
			ReadBufferSize:  socketBufferSize,
			WriteBufferSize: socketBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		room: room,
	}
	handler.room.Run()

	return handler
}

func (h *Handler) generateMessage(sendSummary *baseModels.SendSummary) (result *chat.Message, err error) {
	credentials, err := h.useCase.GetResponseCredentials(sendSummary.SummaryID, sendSummary.VacancyID)
	if err != nil {
		return nil, err
	}

	var status string
	if sendSummary.Accepted {
		status = "одобрено"
	} else if sendSummary.Denied {
		status = "отклонено"
	} else {
		status = "проигнорировано"
	}

	return &chat.Message{
		Message:   fmt.Sprintf("Ваше резюме было %s", status),
		UserOneId: credentials.OrganizationID,
		UserOne:   credentials.OrganizationName,
		UserTwoId: credentials.UserID,
		UserTwo:   credentials.UserName,
		Created:   time.Now(),
	}, nil
}

func (h *Handler) ResponseSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var sendSummary baseModels.SendSummary
	err := json.NewDecoder(r.Body).Decode(&sendSummary)
	if err != nil {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sendSummary.SummaryID, _ = strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	sendSummary.OrganizationID = r.Context().Value("userID").(uint64)

	err = h.useCase.ResponseSummary(&sendSummary)
	switch true {
	case errors.Is(err, summaryInterfaces.ErrOrganizationIsNotOwner):
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusForbidden)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case errors.Is(err, summaryInterfaces.ErrNoSummaryToRefresh):
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		if message, err := h.generateMessage(&sendSummary); err == nil {
			h.room.SendGeneratedMessage(message)
		}
		w.WriteHeader(http.StatusOK)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) EnterChat(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID := r.Context().Value("userID").(uint64)

	socket, err := h.upGrader.Upgrade(w, r, nil)
	if err != nil {
		golog.Errorf("#%s: Failed to start ws - %w", rID, err)
		return
	}

	chatter := &chat.Chatter{
		ID:     userID,
		Socket: socket,
		Send:   make(chan []byte, messageBufferSize),
		Room:   h.room,
	}

	h.room.Join(chatter)
	defer func() {
		h.room.Leave(chatter)
	}()
	go chatter.Write()
	chatter.Read()
}

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	var parameters baseModels.ChatParameters
	rID := r.Context().Value("rID").(string)
	parameters.From = r.Context().Value("userID").(uint64)
	parameters.To, _ = strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)
	parameters.Page, _ = strconv.ParseUint(r.FormValue("page"), 10, 64)

	result, err := h.useCase.GetHistory(&parameters)
	switch true {
	case err == nil:
		w.WriteHeader(http.StatusOK)
		json, _ := json.Marshal(result)
		w.Write(json)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID := r.Context().Value("userID").(uint64)

	result, err := h.useCase.GetConversations(userID)
	switch true {
	case err == nil:
		w.WriteHeader(http.StatusOK)
		json, _ := json.Marshal(result)
		w.Write(json)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}