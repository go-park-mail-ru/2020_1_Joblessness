package httpSummary

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"joblessness/haha/models"
	"joblessness/haha/summary/interfaces"
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
   	var newSummary models.Summary
	newSummary.Author.ID =  r.Context().Value("userID").(uint64)

	err := json.NewDecoder(r.Body).Decode(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newSummary.ID, err = h.useCase.CreateSummary(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(models.ResponseID{ID: newSummary.ID})
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	page := r.FormValue("page")

	summaries, err := h.useCase.GetAllSummaries(page)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) PrintSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)

	getSummary, err := h.useCase.GetSummary(summaryID)
	switch err {
	case sql.ErrNoRows :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		errOut := pdf.SummaryToPdf(w, *getSummary)
		if errOut != nil {
			golog.Errorf("#%s: %w",  rID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/pdf")
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)

	summary, err := h.useCase.GetSummary(summaryID)
	switch err {
	case sql.ErrNoRows :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		jsonData, err := json.Marshal(summary)
		if err != nil {
			golog.Errorf("#%s: %w",  rID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetUserSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, err := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)
	page := r.FormValue("page")

	summaries, err := h.useCase.GetUserSummaries(page, userID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangeSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)

	var newSummary models.Summary
	err := json.NewDecoder(r.Body).Decode(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newSummary.ID = summaryID
	newSummary.Author.ID = r.Context().Value("userID").(uint64)

	err = h.useCase.ChangeSummary(&newSummary)
	switch err {
	case sql.ErrNoRows :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case summaryInterfaces.ErrPersonIsNotOwner:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusForbidden)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	summaryID, _ := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	authorID := r.Context().Value("userID").(uint64)

	err := h.useCase.DeleteSummary(summaryID, authorID)
	switch err {
	case sql.ErrNoRows:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case summaryInterfaces.ErrPersonIsNotOwner:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusForbidden)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) SendSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var sendSummary models.SendSummary
	err := json.NewDecoder(r.Body).Decode(&sendSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sendSummary.VacancyID, _ = strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	sendSummary.UserID =  r.Context().Value("userID").(uint64)

	err = h.useCase.SendSummary(&sendSummary)
	switch err {
	case summaryInterfaces.ErrPersonIsNotOwner:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusForbidden)
	case summaryInterfaces.ErrNoSummaryToRefresh:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ResponseSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var sendSummary models.SendSummary
	err := json.NewDecoder(r.Body).Decode(&sendSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sendSummary.SummaryID, _ = strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	sendSummary.OrganizationID =  r.Context().Value("userID").(uint64)

	err = h.useCase.ResponseSummary(&sendSummary)
	switch err {
	case summaryInterfaces.ErrOrgIsNotOwner:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusForbidden)
	case summaryInterfaces.ErrNoSummaryToRefresh:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetOrgSendSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	summaries, err := h.useCase.GetOrgSendSummaries(userID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetUserSendSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	summaries, err := h.useCase.GetUserSendSummaries(userID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}