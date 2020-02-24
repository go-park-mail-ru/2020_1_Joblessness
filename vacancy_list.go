package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Vacancy struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Skills string `json:"skills"`
	Salary string `json:"salary"`
	Address string `json:"address"`
	PhoneNumber string `json:"phone-number"`
}

type VacancyHandler struct {
	vacancies map[uint]*Vacancy
}

var vacancyId uint

func getNewVacancyId() uint {
	vacancyId++
	return vacancyId
}

func NewVacancyHandler() *VacancyHandler {
	newId := getNewVacancyId()

	return &VacancyHandler {
		vacancies: map[uint]*Vacancy {
			newId: {newId, "name", "description", "skills", "100500", "address", "phone number"},
		},
	}
}

func (api *VacancyHandler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /vacancies")
	Cors.PrivateApi(&w, r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	type Response struct {
		Status uint `json:"status"`
	}

	name := data["name"]
	if name == "" {
		jsonData, _ := json.Marshal(Response{http.StatusBadRequest})
		w.Write(jsonData)
		return
	}

	description := data["description"]
	skills := data["skills"]
	salary := data["salary"]
	address := data["address"]
	phoneNumber := data["phone-number"]

	newId := getNewVacancyId()
	api.vacancies[newId] = &Vacancy{newId, name, description, skills, salary, address, phoneNumber}

	jsonData, _ := json.Marshal(Response{http.StatusCreated})
	w.Write(jsonData)
}

func (api *VacancyHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /vacancies")
	Cors.PrivateApi(&w, r)

	var vacancies []Vacancy
	for _, vacancy := range api.vacancies {
		vacancies = append(vacancies, *vacancy)
	}

	type Response struct {
		Status uint `json:"status"`
		Data []Vacancy `json:"data"`
	}

	jsonData, _ := json.Marshal(Response{http.StatusOK, vacancies})
	w.Write(jsonData)
}

func (api *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	type Response struct {
		Status uint `json:"status"`
		Data Vacancy `json:"data,omitempty"`
	}

	vacancy, ok := api.vacancies[uint(vacancyId)]
	if !ok {
		jsonData, _ := json.Marshal(Response{Status:http.StatusNotFound})
		w.Write(jsonData)
		return
	}

	jsonData, _ := json.Marshal(Response{http.StatusOK, *vacancy})
	w.Write(jsonData)
}

func (api *VacancyHandler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	type Response struct {
		Status uint `json:"status"`
	}

	if _, ok := api.vacancies[uint(vacancyId)]; !ok {
		jsonData, _ := json.Marshal(Response{http.StatusNotFound})
		w.Write(jsonData)
		return
	}

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	name := data["name"]
	if name == "" {
		jsonData, _ := json.Marshal(Response{Status:http.StatusBadRequest})
		w.Write(jsonData)
		return
	}

	description := data["description"]
	skills := data["skills"]
	salary := data["salary"]
	address := data["address"]
	phoneNumber := data["phone-number"]

	api.vacancies[uint(vacancyId)] = &Vacancy{uint(vacancyId), name, description, skills, salary, address, phoneNumber}

	jsonData, _ := json.Marshal(Response{http.StatusNoContent})
	w.Write(jsonData)
}

func (api *VacancyHandler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE /vacancies/{vacancy_id}")
	Cors.PrivateApi(&w, r)

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	type Response struct {
		Status uint `json:"status"`
	}

	if _, ok := api.vacancies[uint(vacancyId)]; !ok {
		jsonData, _ := json.Marshal(Response{http.StatusNotFound})
		w.Write(jsonData)
	}

	delete(api.vacancies, uint(vacancyId))

	jsonData, _ := json.Marshal(Response{http.StatusOK})
	w.Write(jsonData)
}