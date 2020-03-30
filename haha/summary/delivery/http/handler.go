package httpSummary

import (
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

type Response struct {
	ID uint64 `json:"id"`
}

func (h *Handler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var newSummary models.Summary
	err := json.NewDecoder(r.Body).Decode(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO проверять существование контекста
	newSummary.Author.ID =  r.Context().Value("userID").(uint64)

	summaryID, err := h.useCase.CreateSummary(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(Response{summaryID})
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

	summaries, err := h.useCase.GetAllSummaries()
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

func (h *Handler) PrintSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	getSummary, err := h.useCase.GetSummary(summaryID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	errOut := pdf.SummaryToPdf(w, *getSummary)
	if errOut != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/pdf")
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summary, err := h.useCase.GetSummary(summaryID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	golog.Infof("#%s: %w", rID, summary)

	jsonData, err := json.Marshal(summary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetUserSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, err := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summaries, err := h.useCase.GetUserSummaries(userID)
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

	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newSummary models.Summary
	err = json.NewDecoder(r.Body).Decode(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newSummary.ID = summaryID
	newSummary.Author.ID = r.Context().Value("userID").(uint64)

	err = h.useCase.ChangeSummary(&newSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.DeleteSummary(summaryID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SendSummary(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var sendSummary models.SendSummary
	err := json.NewDecoder(r.Body).Decode(&sendSummary)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendSummary.VacancyID, err = strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendSummary.SummaryID, err = strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

func (h *Handler) GetOrgSummaries(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, err := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summaries, err := h.useCase.GetOrgSummaries(userID)
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