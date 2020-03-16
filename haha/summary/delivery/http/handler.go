package httpSummary

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"joblessness/haha/models"
	"joblessness/haha/summary"
	"joblessness/haha/utils/pdf"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase summary.UseCase
}

func NewHandler(useCase summary.UseCase) *Handler {
	return &Handler{useCase}
}

type Response struct {
	ID uint64 `json:"id"`
}

func (h *Handler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	summary := models.Summary{
		UserID: r.Context().Value("userID").(uint64),
	}

	err = json.Unmarshal(body, &summary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	summaryID, err := h.useCase.CreateSummary(&summary)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(Response{summaryID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	summaries, err := h.useCase.GetAllSummaries()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) PrintSummary(w http.ResponseWriter, r *http.Request) {
	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summary, err := h.useCase.GetSummary(summaryID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	errOut := pdf.SummaryToPdf(w, *summary)
	if errOut != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/pdf")
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summary, err := h.useCase.GetSummary(summaryID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(summary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetUserSummaries(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summaries, err := h.useCase.GetUSerSummaries(userID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangeSummary(w http.ResponseWriter, r *http.Request) {
	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	summary := models.Summary{
		ID: summaryID,
		UserID: r.Context().Value("userID").(uint64),
	}

	err = json.Unmarshal(body, &summary)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.useCase.ChangeSummary(&summary)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	summaryID, err := strconv.ParseUint(mux.Vars(r)["summary_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.DeleteSummary(summaryID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
