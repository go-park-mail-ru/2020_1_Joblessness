package summaries

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func NewEmptySummaryHandler() *SummaryHandler {
	return &SummaryHandler{
		summaries:map[uint]*Summary{},
		mu: sync.RWMutex{},
	}
}

func NewNotEmptySummaryHandler() *SummaryHandler {
	return &SummaryHandler{
		summaries:map[uint]*Summary{
			1: {1, 1, "first name", "last name", "phone number", "email", "birth date", "male", "experiencz", "education"},
			2: {1, 2, "first name", "last name", "phone number", "email", "birth date", "female", "experiencz", "education"},
		},
		mu: sync.RWMutex{},
		SummaryId:2,
	}
}

func TestCreateSummaryFailedNoAuthor(t *testing.T) {
	t.Parallel()

	h := NewEmptySummaryHandler()

	summary, _ := json.Marshal(Summary{
		FirstName:   "first name",
		LastName:    "last name",
		PhoneNumber: "phone number",
		Email:       "email",
		BirthDate:   "birth date",
		Gender:      "gender",
		Experience:  "experience",
		Education:   "education",
	})

	body := bytes.NewReader([]byte(summary))

	r := httptest.NewRequest("POST", "/api/summaries", body)
	w := httptest.NewRecorder()

	h.CreateSummary(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	if len(h.summaries) != 0 {
		t.Error("Wrong Summary created")
	}
	//
	//if h.summaries[1].FirstName != "first name" {
	//	t.Error("Wrong summary first name")
	//}
}

func TestGetEmptySummaryList(t *testing.T) {
	t.Parallel()

	h := NewEmptySummaryHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/summaries", body)
	w := httptest.NewRecorder()

	h.GetSummaries(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not 200")
	}
}

func TestGetNotEmptySummaryList(t *testing.T) {
	t.Parallel()

	h := NewNotEmptySummaryHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/summaries", body)
	w := httptest.NewRecorder()

	h.GetSummaries(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not 200")
	}

	if !strings.Contains(w.Body.String(), "male") {
		t.Error("First summary is not in list")
	}

	if !strings.Contains(w.Body.String(), "female") {
		t.Error("Second summary is not in list")
	}
}

func TestSuccessGetSummary(t *testing.T) {
	t.Parallel()

	h := NewNotEmptySummaryHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/summaries/1", body)
	r = mux.SetURLVars(r, map[string]string{"summary_id": "1"})
	w := httptest.NewRecorder()

	h.GetSummary(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not 200")
	}

	if !strings.Contains(w.Body.String(), "male") {
		t.Error("First summary is not in list")
	}

	if strings.Contains(w.Body.String(), "female") {
		t.Error("Second summary is in list")
	}
}

func TestGetNullSummary(t *testing.T) {
	t.Parallel()

	h := NewNotEmptySummaryHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/summaries/3", body)
	r = mux.SetURLVars(r, map[string]string{"summary_id": "3"})
	w := httptest.NewRecorder()

	h.GetSummary(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}
}

func TestSuccessChangeSummary(t *testing.T) {
	t.Parallel()

	h := NewNotEmptySummaryHandler()

	summary, _ := json.Marshal(Summary{
		FirstName:   "first name",
		LastName:    "last name",
		PhoneNumber: "phone number",
		Email:       "email",
		BirthDate:   "birth date",
		Gender:      "new gender",
		Experience:  "experience",
		Education:   "education",
	})

	body := bytes.NewReader([]byte(summary))

	r := httptest.NewRequest("PUT", "/api/summaries/1", body)
	r = mux.SetURLVars(r, map[string]string{"summary_id": "1"})
	w := httptest.NewRecorder()

	h.ChangeSummary(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	if h.summaries[1].Gender != "male" {
		t.Error("Wrong Vacancy is changed")
	}
}

func TestChangeNullSummary(t *testing.T) {
	t.Parallel()

	h := NewNotEmptySummaryHandler()

	summary, _ := json.Marshal(Summary{
		FirstName:   "first name",
		LastName:    "last name",
		PhoneNumber: "phone number",
		Email:       "email",
		BirthDate:   "birth date",
		Gender:      "new gender",
		Experience:  "experience",
		Education:   "education",
	})

	body := bytes.NewReader([]byte(summary))

	r := httptest.NewRequest("PUT", "/api/summaries/3", body)
	r = mux.SetURLVars(r, map[string]string{"summary_id": "3"})
	w := httptest.NewRecorder()

	h.ChangeSummary(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}
}

func TestSuccessDeleteSummary(t *testing.T) {
	t.Parallel()

	h := NewNotEmptySummaryHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("DELETE", "/api/summaries/1", body)
	r = mux.SetURLVars(r, map[string]string{"summary_id": "1"})
	w := httptest.NewRecorder()

	h.DeleteSummary(w, r)

	if w.Code != http.StatusNoContent {
		t.Error("Status code is not 204")
	}

	if len(h.summaries) != 1 {
		t.Error("Summary is not deleted")
	}
}

func TestDeleteNullSummary(t *testing.T) {
	t.Parallel()

	h := NewNotEmptySummaryHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("DELETE", "/api/summaries/3", body)
	r = mux.SetURLVars(r, map[string]string{"summary_id": "3"})
	w := httptest.NewRecorder()

	h.DeleteSummary(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status code is not 404")
	}
}