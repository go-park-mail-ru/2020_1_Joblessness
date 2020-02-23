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

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	name := data["name"]
	if name == "" {
		http.Error(w, `Empty name`, 400)
		return
	}

	description := data["description"]
	skills := data["skills"]
	salary := data["salary"]
	address := data["address"]
	phoneNumber := data["phone-number"]

	newId := getNewVacancyId()
	api.vacancies[newId] = &Vacancy{newId, name, description, skills, salary, address, phoneNumber}
}

func (api *VacancyHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /vacancies")

	var vacancies []Vacancy
	for _, vacancy := range api.vacancies {
		vacancies = append(vacancies, *vacancy)
	}

	jsonData, _ := json.Marshal(vacancies)

	w.Write(jsonData)
}

func (api *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /vacancies/{vacancy_id}")

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	vacancy, ok := api.vacancies[uint(vacancyId)]
	if !ok {
		http.Error(w, `Not found`, 404)
		return
	}

	jsonData, _ := json.Marshal(vacancy)

	w.Write(jsonData)
}

func (api *VacancyHandler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /vacancies/{vacancy_id}")

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	if _, ok := api.vacancies[uint(vacancyId)]; !ok {
		http.Error(w, `Not found`, 404)
		return
	}

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	name := data["name"]
	if name == "" {
		http.Error(w, `Empty name`, 400)
		return
	}

	description := data["description"]
	skills := data["skills"]
	salary := data["salary"]
	address := data["address"]
	phoneNumber := data["phone-number"]

	api.vacancies[uint(vacancyId)] = &Vacancy{uint(vacancyId), name, description, skills, salary, address, phoneNumber}
}

func (api *VacancyHandler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE /vacancies/{vacancy_id}")

	vacancyId, _ := strconv.Atoi(mux.Vars(r)["vacancy_id"])

	delete(api.vacancies, uint(vacancyId))
}