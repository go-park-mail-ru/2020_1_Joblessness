package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Summary struct {
	ID uint `json:"id,omitempty"`
	FirstName string `json:"first-name"`
	LastName string `json:"last-name"`
	PhoneNumber string `json:"phone-number"`
	Email string `json:"email"`
	BirthDate string `json:"birth-date"`
	Gender string `json:"gender"`
	Experience string `json:"experience"`
	Education string `json:"education"`
	Author uint `json:"author,omitempty"`
}

type SummaryHandler struct {
	summaries map[uint]*Summary
	SummaryId uint
}

func (api *SummaryHandler) getNewSummaryId() uint {
	api.SummaryId++
	return api.SummaryId
}

func NewSummaryHandler() *SummaryHandler {
	return &SummaryHandler {
		summaries: map[uint]*Summary {
			1: {1, "first name", "last name", "phone number", "kek@mail.ru", "01/01/1900", "gender", "experience", "bmstu", 1},
		},
		SummaryId:1,
	}
}

func (api *SummaryHandler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /summaries")
	Cors.PrivateApi(&w, r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	log.Println("recieved: ", data)
	author, found := data["userid"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authorId, err := strconv.Atoi(author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newId := api.getNewSummaryId()
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
		uint(authorId),
	}

	type Response struct {
		ID uint `json:"id"`
	}

	jsonData, _ := json.Marshal(Response{newId})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (api *SummaryHandler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /summaries")
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
	log.Println("GET /summaries/{summary_id}")
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

func (api *SummaryHandler) GetUserSummaries(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /users/{user_id}/summaries")
	Cors.PrivateApi(&w, r)

	userId, _ := strconv.Atoi(mux.Vars(r)["user_id"])

	var summaries []Summary
	for _, summary := range api.summaries {
		if (*summary).Author == uint(userId) {
			summaries = append(summaries, *summary)
		}
	}

	if len(summaries) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonData, _ := json.Marshal(summaries)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *SummaryHandler) ChangeSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, _ := strconv.Atoi(mux.Vars(r)["summary_id"])

	if _, ok := api.summaries[uint(summaryId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	author, found := data["userid"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authorId, err := strconv.Atoi(author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
		uint(authorId),
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *SummaryHandler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, _ := strconv.Atoi(mux.Vars(r)["summary_id"])

	if _, ok := api.summaries[uint(summaryId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(api.summaries, uint(summaryId))

	w.WriteHeader(http.StatusNoContent)
}