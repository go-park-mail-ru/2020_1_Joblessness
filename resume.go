package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

type Summary struct {
	UserID uint `json:"author,omitempty"`
	ID uint `json:"id,omitempty"`
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
	mu sync.RWMutex
	SummaryId uint32
}

func (api *SummaryHandler) getNewSummaryId() uint32 {
	return atomic.AddUint32(&api.SummaryId, 1)
}

func NewSummaryHandler() *SummaryHandler {
	return &SummaryHandler {
		summaries: map[uint]*Summary {
			1: {1, 1, "first name", "last name", "phone number", "kek@mail.ru", "01/01/1900", "gender", "experience", "bmstu"},
		},
		mu: sync.RWMutex{},
		SummaryId:1,
	}
}

func (api *SummaryHandler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /summaries")
	Cors.PrivateApi(&w, r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	log.Println(data)
	author, found := data["author"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authorId, err := strconv.Atoi(author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.Lock()
	newId := api.getNewSummaryId()
	api.summaries[uint(newId)] = &Summary{
		uint(authorId),
		uint(newId),
		data["first-name"],
		data["last-name"],
		data["phone-number"],
		data["email"],
		data["birth-date"],
		data["gender"],
		data["experience"],
		data["education"],
	}
	api.mu.Unlock()

	type Response struct {
		ID uint32 `json:"id"`
	}

	jsonData, _ := json.Marshal(Response{newId})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (api *SummaryHandler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /summaries")
	Cors.PrivateApi(&w, r)

	var summaries []Summary
	api.mu.RLock()
	for _, summary := range api.summaries {
		summaries = append(summaries, *summary)
	}
	api.mu.RUnlock()

	jsonData, _ := json.Marshal(summaries)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *SummaryHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, _ := strconv.Atoi(mux.Vars(r)["summary_id"])

	api.mu.RLock()
	summary, ok := api.summaries[uint(summaryId)]
	api.mu.RUnlock()
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
	api.mu.RLock()
	for _, summary := range api.summaries {
		if (*summary).UserID == uint(userId) {
			summaries = append(summaries, *summary)
		}
	}
	api.mu.RUnlock()

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

	author, found := data["author"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authorId, err := strconv.Atoi(author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.Lock()
	api.summaries[uint(summaryId)] = &Summary{
		uint(authorId),
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
	api.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func (api *SummaryHandler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, _ := strconv.Atoi(mux.Vars(r)["summary_id"])

	api.mu.Lock()
	if _, ok := api.summaries[uint(summaryId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(api.summaries, uint(summaryId))
	api.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}