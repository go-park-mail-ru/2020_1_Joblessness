package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"joblessness/haha/models"
	"joblessness/haha/vacancy"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase vacancy.UseCase
}

func NewHandler(useCase vacancy.UseCase) *Handler {
	return &Handler{useCase}
}

type Response struct {
	ID uint64 `json:"id"`
}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	var vacancy models.Vacancy
	json.NewDecoder(r.Body).Decode(&vacancy)

	if vacancy.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vacancyID, err := h.useCase.CreateVacancy(vacancy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(Response{vacancyID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, err := h.useCase.GetVacancies()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(vacancies) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonData, err := json.Marshal(vacancies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	vacancyId, err := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vacancy, err := h.useCase.GetVacancy(vacancyId)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(vacancy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /vacancies/{vacancy_id}")

	vacancyID, err := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
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

	var vacancy models.Vacancy

	err = json.Unmarshal(body, &vacancy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if vacancy.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vacancy.ID = vacancyID

	err = h.useCase.ChangeVacancy(vacancy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /vacancies/{vacancy_id}")

	vacancyID, err := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.DeleteVacancy(vacancyID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
