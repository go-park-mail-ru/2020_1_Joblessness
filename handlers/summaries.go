package handlers

import (
	_models "../models"
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
	Summaries map[uint]*_models.Summary
	Mu        sync.RWMutex
	SummaryId uint32
}

func (api *SummaryHandler) getNewSummaryId() uint32 {
	return atomic.AddUint32(&api.SummaryId, 1)
}

func NewSummaryHandler() *SummaryHandler {
	return &SummaryHandler {
		Summaries: map[uint]*_models.Summary {
			1: {1, 1, "first name", "last name", "phone number", "kek@mail.ru", "01/01/1900", "gender", "experience", "bmstu"},
		},
		Mu:        sync.RWMutex{},
		SummaryId: 1,
	}
}

func (api *SummaryHandler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /summaries")

	newId := api.getNewSummaryId()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var summary _models.Summary
	err = json.Unmarshal(body, &summary)
	summary.ID = uint(newId)
	log.Println("summary recieved: ", summary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	api.Mu.Lock()
	api.Summaries[uint(newId)] = &summary
	api.Mu.Unlock()

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

	var summaries []_models.Summary
	api.Mu.RLock()
	for _, summary := range api.Summaries {
		summaries = append(summaries, *summary)
	}
	api.Mu.RUnlock()

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

	summaryId, err := strconv.Atoi(mux.Vars(r)["summary_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.Mu.RLock()
	summary, ok := api.Summaries[uint(summaryId)]
	api.Mu.RUnlock()
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

	userId, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var summaries []_models.Summary
	api.Mu.RLock()
	for _, summary := range api.Summaries {
		if (*summary).UserID == uint(userId) {
			summaries = append(summaries, *summary)
		}
	}
	api.Mu.RUnlock()

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

	summaryId, err := strconv.Atoi(mux.Vars(r)["summary_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := api.Summaries[uint(summaryId)]; !ok {
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

	api.Mu.Lock()
	api.Summaries[uint(summaryId)] = &_models.Summary{
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
	api.Mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func (api *SummaryHandler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /summaries/{summary_id}")

	summaryId, err := strconv.Atoi(mux.Vars(r)["summary_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.Mu.Lock()
	if _, ok := api.Summaries[uint(summaryId)]; !ok {
		w.WriteHeader(http.StatusNotFound)
		api.Mu.Unlock()
		return
	}

	delete(api.Summaries, uint(summaryId))
	api.Mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}