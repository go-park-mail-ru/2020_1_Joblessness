package handlers

import (
	_models "../models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

type VacancyHandler struct {
	vacancies map[uint]*_models.Vacancy
	mu sync.RWMutex
	vacancyId uint32
}

func (api *VacancyHandler) getNewVacancyId() uint32 {
	return atomic.AddUint32(&api.vacancyId, 1)
}

func NewVacancyHandler() *VacancyHandler {
	return &VacancyHandler {
		vacancies: map[uint]*_models.Vacancy {
			1: {1, "name", "description", "skills", "100500", "address", "phone number"},
		},
		mu: sync.RWMutex{},
		vacancyId:1,
	}
}

func (api *VacancyHandler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /vacancies")
	Cors.PrivateApi(&w, r)

	var vacancy _models.Vacancy
	json.NewDecoder(r.Body).Decode(&vacancy)

	if vacancy.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newId := api.getNewVacancyId()

	api.mu.Lock()
	api.vacancies[uint(newId)] = &vacancy
	api.mu.Unlock()

	type Response struct {
		ID uint32 `json:"id"`
	}

	jsonData, err := json.Marshal(Response{newId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (api *VacancyHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /vacancies")
	Cors.PrivateApi(&w, r)

	var vacancies []_models.Vacancy
	api.mu.RLock()
	for _, vacancy := range api.vacancies {
		vacancies = append(vacancies, *vacancy)
	}
	api.mu.RUnlock()

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

func (api *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, err := strconv.Atoi(mux.Vars(r)["vacancy_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.RLock()
	vacancy, ok := api.vacancies[uint(vacancyId)]
	api.mu.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
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

func (api *VacancyHandler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, err := strconv.Atoi(mux.Vars(r)["vacancy_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := api.vacancies[uint(vacancyId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	name := data["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	description := data["description"]
	skills := data["skills"]
	salary := data["salary"]
	address := data["address"]
	phoneNumber := data["phone-number"]

	api.mu.Lock()
	api.vacancies[uint(vacancyId)] = &_models.Vacancy{uint(vacancyId), name, description, skills, salary, address, phoneNumber}
	api.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func (api *VacancyHandler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, err := strconv.Atoi(mux.Vars(r)["vacancy_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.Lock()
	if _, ok := api.vacancies[uint(vacancyId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(api.vacancies, uint(vacancyId))
	api.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
