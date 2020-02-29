package main

import (
	_models "./models"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

type SummaryHandler struct {
	summaries map[uint]*_models.Summary
	mu sync.RWMutex
	SummaryId uint32
}

func (api *SummaryHandler) getNewSummaryId() uint32 {
	return atomic.AddUint32(&api.SummaryId, 1)
}

func NewSummaryHandler() *SummaryHandler {
	return &SummaryHandler {
		summaries: map[uint]*_models.Summary {
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
	_, err := strconv.Atoi(author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newId := api.getNewSummaryId()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var summary _models.Summary
	err = json.Unmarshal(body, &summary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	api.mu.Lock()
	api.summaries[uint(newId)] = &summary
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

func (api *SummaryHandler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /summaries")
	Cors.PrivateApi(&w, r)

	var summaries []_models.Summary
	api.mu.RLock()
	for _, summary := range api.summaries {
		summaries = append(summaries, *summary)
	}
	api.mu.RUnlock()

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *SummaryHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, err := strconv.Atoi(mux.Vars(r)["summary_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.RLock()
	summary, ok := api.summaries[uint(summaryId)]
	api.mu.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
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

func (api *SummaryHandler) GetUserSummaries(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /users/{user_id}/summaries")
	Cors.PrivateApi(&w, r)

	userId, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var summaries []_models.Summary
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

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *SummaryHandler) ChangeSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /summaries/{summary_id}")
	Cors.PrivateApi(&w, r)

	summaryId, err := strconv.Atoi(mux.Vars(r)["summary_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	api.summaries[uint(summaryId)] = &_models.Summary{
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

	summaryId, err := strconv.Atoi(mux.Vars(r)["summary_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.Lock()
	if _, ok := api.summaries[uint(summaryId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(api.summaries, uint(summaryId))
	api.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}