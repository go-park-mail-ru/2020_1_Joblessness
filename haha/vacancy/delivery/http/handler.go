package httpVacancy

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"joblessness/haha/models"
	"joblessness/haha/vacancy"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase vacancy.VacancyUseCase
}

func NewHandler(useCase vacancy.VacancyUseCase) *Handler {
	return &Handler{useCase}
}

type Response struct {
	ID uint64 `json:"id"`
}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var newVacancy models.Vacancy
	err := json.NewDecoder(r.Body).Decode(&newVacancy)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if newVacancy.Name == "" {
		golog.Errorf("#%s: %s",  rID, "empty vacancy name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO проверять существование контекста
	newVacancy.Organization.ID =  r.Context().Value("userID").(uint64)

	vacancyID, err := h.useCase.CreateVacancy(&newVacancy)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(Response{vacancyID})
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	vacancyId, err := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	getVacancy, err := h.useCase.GetVacancy(vacancyId)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(getVacancy)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	vacancies, err := h.useCase.GetVacancies()
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(vacancies) == 0 {
		golog.Errorf("#%s: %s",  rID, "no vacancies")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonData, err := json.Marshal(vacancies)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	vacancyID, err := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newVacancy models.Vacancy
	err = json.NewDecoder(r.Body).Decode(&newVacancy)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if newVacancy.Name == "" {
		golog.Errorf("#%s: %s",  rID, "empty vacancy name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newVacancy.ID = vacancyID

	err = h.useCase.ChangeVacancy(&newVacancy)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	vacancyID, err := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.DeleteVacancy(vacancyID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
