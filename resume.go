package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Summary struct {
	ID uint `json:"id"`
	FirstName string `json:"first-name"`
	LastName string `json:"last-name"`
	PhoneNumber string `json:"phone-number"`
	Email string `json:"email"`
	BirthDate string `json:"birth-date"`
	Gender string `json:"gender"`
	Experience string `json:"experience"`
	Education string `json:"education"`
}

type SummaryHandler struct {
	summaries map[uint]*Summary
}

var SummaryId uint

func getNewSummaryId() uint {
	SummaryId++
	return SummaryId
}

func NewSummaryHandler() *SummaryHandler {
	newId := getNewSummaryId()

	return &SummaryHandler {
		summaries: map[uint]*Summary {
			newId: {newId, "first name", "last name", "phone number", "kek@mail.ru", "01/01/1900", "gender", "experience", "bmstu"},
		},
	}
}

func (api *SummaryHandler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /summaries")
	Cors.PrivateApi(&w, r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	newId := getNewSummaryId()
	api.summaries[newId] = &Summary{
		newId,
		data["first-name"],
		data["last-name"],
		data["phone-number"],
		data["email"],
		data["birth-date"],
		data["gender"],
		data["experience"],
		data["education"],
	}

	w.WriteHeader(http.StatusCreated)
}

func (api *SummaryHandler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /summaries")
	Cors.PrivateApi(&w, r)

	var summaries []Summary
	for _, summary := range api.summaries {
		summaries = append(summaries, *summary)
	}

	jsonData, _ := json.Marshal(summaries)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *SummaryHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, _ := strconv.Atoi(mux.Vars(r)["summary_id"])

	summary, ok := api.summaries[uint(summaryId)]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonData, _ := json.Marshal(summary)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *SummaryHandler) ChangeSummary(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, _ := strconv.Atoi(mux.Vars(r)["summary_id"])

	if _, ok := api.summaries[uint(summaryId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	api.summaries[uint(summaryId)] = &Summary{
		uint(summaryId),
		data["first-name"],
		data["last-name"],
		data["phone-number"],
		data["email"],
		data["birth-date"],
		data["gender"],
		data["experience"],
		data["education"],
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *SummaryHandler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, _ := strconv.Atoi(mux.Vars(r)["summary_id"])

	if _, ok := api.summaries[uint(summaryId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(api.summaries, uint(summaryId))

	w.WriteHeader(http.StatusNoContent)
}