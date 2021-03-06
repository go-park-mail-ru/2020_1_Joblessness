package summaryHttp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/mailru/easyjson"
	"joblessness/haha/models/base"
	"joblessness/haha/summary/interfaces"
	"joblessness/haha/utils/mail"
	"joblessness/haha/utils/pdf"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase summaryInterfaces.SummaryUseCase
}

func NewHandler(useCase summaryInterfaces.SummaryUseCase) *Handler {
	return &Handler{useCase}
}

func (h *Handler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	var newSummary baseModels.Summary
	newSummary.Author.ID = r.Context().Value("userID").(uint64)

	err := easyjson.UnmarshalFromReader(r.Body, &newSummary)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newSummary.ID, err = h.useCase.CreateSummary(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(baseModels.ResponseID{ID: newSummary.ID})

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	page := r.FormValue("page")

	summaries, err := h.useCase.GetAllSummaries(page)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(summaries)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) PrintSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)

	getSummary, err := h.useCase.GetSummary(summaryID)
	switch true {
	case errors.Is(err, sql.ErrNoRows):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusNotFound)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		w.Header().Set("Content-Type", "application/pdf")
		errOut := pdf.SummaryToPdf(w, getSummary)
		if errOut != nil {
			golog.Errorf("#%s: %s", rID, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)

	summary, err := h.useCase.GetSummary(summaryID)
	switch true {
	case errors.Is(err, sql.ErrNoRows):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusNotFound)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		jsonData, _ := easyjson.Marshal(summary)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) GetUserSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, err := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := r.FormValue("page")

	summaries, err := h.useCase.GetUserSummaries(page, userID)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(summaries)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangeSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)

	var newSummary baseModels.Summary
	err := easyjson.UnmarshalFromReader(r.Body, &newSummary)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newSummary.ID = summaryID
	newSummary.Author.ID = r.Context().Value("userID").(uint64)

	err = h.useCase.ChangeSummary(&newSummary)
	switch true {
	case errors.Is(err, summaryInterfaces.ErrSummaryNotFound):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusNotFound)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case errors.Is(err, summaryInterfaces.ErrPersonIsNotOwner):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusForbidden)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	authorID := r.Context().Value("userID").(uint64)

	err := h.useCase.DeleteSummary(summaryID, authorID)
	switch true {
	case errors.Is(err, summaryInterfaces.ErrSummaryNotFound):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusNotFound)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case errors.Is(err, summaryInterfaces.ErrPersonIsNotOwner):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusForbidden)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) SendSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var sendSummary baseModels.SendSummary
	err := easyjson.UnmarshalFromReader(r.Body, &sendSummary)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sendSummary.VacancyID, _ = strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	sendSummary.UserID = r.Context().Value("userID").(uint64)

	err = h.useCase.SendSummary(&sendSummary)
	switch true {
	case errors.Is(err, summaryInterfaces.ErrPersonIsNotOwner):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusForbidden)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case errors.Is(err, summaryInterfaces.ErrNoSummaryToRefresh):
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusNotFound)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		w.WriteHeader(http.StatusOK)
	default:
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := easyjson.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) GetOrgSendSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	summaries, err := h.useCase.GetOrgSendSummaries(userID)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(summaries)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetUserSendSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	summaries, err := h.useCase.GetUserSendSummaries(userID)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(summaries)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) SendSummaryByMail(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	authorID := r.Context().Value("userID").(uint64)

	var mail mail.Mail
	err := json.NewDecoder(r.Body).Decode(&mail)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.SendSummaryByMail(summaryID, authorID, mail.To)
	switch true {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	default:
		golog.Errorf("#%s: %s", rID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}
