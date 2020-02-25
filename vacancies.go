package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Vacancy struct {
	ID uint `json:"id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description"`
	Skills string `json:"skills"`
	Salary string `json:"salary"`
	Address string `json:"address"`
	PhoneNumber string `json:"phone-number"`
}

type VacancyHandler struct {
	vacancies map[uint]*Vacancy
	vacancyId uint
}

func (api *VacancyHandler) getNewVacancyId() uint {
	api.vacancyId++
	return api.vacancyId
}

func NewVacancyHandler() *VacancyHandler {
	return &VacancyHandler {
		vacancies: map[uint]*Vacancy {
			1: {1, "name", "description", "skills", "100500", "address", "phone number"},
		},
		vacancyId:1,
	}
}

func (api *VacancyHandler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /vacancies")
	Cors.PrivateApi(&w, r)

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

	newId := api.getNewVacancyId()
	api.vacancies[newId] = &Vacancy{newId, name, description, skills, salary, address, phoneNumber}

	w.WriteHeader(http.StatusCreated)
}

func (api *VacancyHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /vacancies")
	Cors.PrivateApi(&w, r)

	var vacancies []Vacancy
	for _, vacancy := range api.vacancies {
		vacancies = append(vacancies, *vacancy)
	}

	if len(vacancies) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonData, _ := json.Marshal(vacancies)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	vacancy, ok := api.vacancies[uint(vacancyId)]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonData, _ := json.Marshal(vacancy)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *VacancyHandler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

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

	api.vacancies[uint(vacancyId)] = &Vacancy{uint(vacancyId), name, description, skills, salary, address, phoneNumber}

	w.WriteHeader(http.StatusNoContent)
}

func (api *VacancyHandler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	if _, ok := api.vacancies[uint(vacancyId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(api.vacancies, uint(vacancyId))

	w.WriteHeader(http.StatusNoContent)
}
