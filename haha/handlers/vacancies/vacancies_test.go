package vacancies

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

func NewEmptyVacancyHandler() *VacancyHandler {
	return &VacancyHandler{
		vacancies: map[uint]*Vacancy{},
		mu: sync.RWMutex{},
	}
}

func NewNotEmptyVacancyHandler() *VacancyHandler {
	return &VacancyHandler{
		vacancies:map[uint]*Vacancy{
			1: {1, "first name", "description", "skills", "salary", "address", "phone"},
			2: {2, "second name", "description", "skills", "salary", "address", "phone"},
		},
		mu: sync.RWMutex{},
		vacancyId:2,
	}
}

func TestSuccessCreateVacancy(t *testing.T) {
	t.Parallel()

	h := NewEmptyVacancyHandler()

	vacancy, _ := json.Marshal(Vacancy{
		Name:        "name",
		Description: "description",
		Skills:      "skills",
		Salary:      "salary",
		Address:     "address",
		PhoneNumber: "phone number",
	})

	body := bytes.NewReader([]byte(vacancy))

	r := httptest.NewRequest("POST", "/api/vacancies", body)
	w := httptest.NewRecorder()

	h.CreateVacancy(w, r)

	if w.Code != http.StatusCreated {
		t.Error("Status is not 201")
	}

	if len(h.vacancies) != 1 {
		t.Error("Vacancy not created")
	}

	if h.vacancies[1].Name != "name" {
		t.Error("Wrong vacancy name")
	}
}

func TestCreateVacancyWithEmptyName(t *testing.T) {
	t.Parallel()

	h := NewEmptyVacancyHandler()

	vacancy, _ := json.Marshal(Vacancy{
		Name:        "",
		Description: "description",
		Skills:      "skills",
		Salary:      "salary",
		Address:     "address",
		PhoneNumber: "phone number",
	})

	body := bytes.NewReader([]byte(vacancy))

	r := httptest.NewRequest("POST", "/api/vacancies", body)
	w := httptest.NewRecorder()

	h.CreateVacancy(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	if len(h.vacancies) != 0 {
		t.Error("Vacancy created")
	}
}

func TestGetEmptyVacancyList(t *testing.T) {
	t.Parallel()

	h := NewEmptyVacancyHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/vacancies", body)
	w := httptest.NewRecorder()

	h.GetVacancies(w, r)

	if w.Code != http.StatusNoContent {
		t.Error("Status is not 204")
	}
}

func TestGetNotEmptyVacancyList(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/vacancies", body)
	w := httptest.NewRecorder()

	h.GetVacancies(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not 200")
	}

	if !strings.Contains(w.Body.String(), "first name") {
		t.Error("First vacancy is not in list")
	}

	if !strings.Contains(w.Body.String(), "second name") {
		t.Error("Second vacancy is not in list")
	}
}

func TestSuccessGetVacancy(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/vacancies/1", body)
	r = mux.SetURLVars(r, map[string]string{"vacancy_id": "1"})
	w := httptest.NewRecorder()

	h.GetVacancy(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not 200")
	}

	if !strings.Contains(w.Body.String(), "first name") {
		t.Error("First vacancy is not in list")
	}

	if strings.Contains(w.Body.String(), "second name") {
		t.Error("Second vacancy is in list")
	}
}

func TestGetNullVacancy(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("GET", "/api/vacancies/3", body)
	r = mux.SetURLVars(r, map[string]string{"vacancy_id": "3"})
	w := httptest.NewRecorder()

	h.GetVacancy(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}
}

func TestSuccessChangeVacancy(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	vacancy, _ := json.Marshal(Vacancy{
		Name:        "new name",
		Description: "description",
		Skills:      "skills",
		Salary:      "salary",
		Address:     "address",
		PhoneNumber: "phone number",
	})

	body := bytes.NewReader([]byte(vacancy))

	r := httptest.NewRequest("PUT", "/api/vacancies/1", body)
	r = mux.SetURLVars(r, map[string]string{"vacancy_id": "1"})
	w := httptest.NewRecorder()

	h.ChangeVacancy(w, r)

	if w.Code != http.StatusNoContent {
		t.Error("Status is not 204")
	}

	if h.vacancies[1].Name != "new name" {
		t.Error("Vacancy is not changed")
	}
}

func TestChangeVacancyWithEmptyName(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	vacancy, _ := json.Marshal(Vacancy{
		Name:        "",
		Description: "description",
		Skills:      "skills",
		Salary:      "salary",
		Address:     "address",
		PhoneNumber: "phone number",
	})

	body := bytes.NewReader([]byte(vacancy))

	r := httptest.NewRequest("PUT", "/api/vacancies/1", body)
	r = mux.SetURLVars(r, map[string]string{"vacancy_id": "1"})
	w := httptest.NewRecorder()

	h.ChangeVacancy(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	if h.vacancies[1].Name != "first name" {
		t.Error("Vacancy is changed")
	}
}

func TestChangeNullVacancy(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	vacancy, _ := json.Marshal(Vacancy{
		Name:        "new name",
		Description: "description",
		Skills:      "skills",
		Salary:      "salary",
		Address:     "address",
		PhoneNumber: "phone number",
	})

	body := bytes.NewReader([]byte(vacancy))

	r := httptest.NewRequest("PUT", "/api/vacancies/3", body)
	r = mux.SetURLVars(r, map[string]string{"vacancy_id": "3"})
	w := httptest.NewRecorder()

	h.ChangeVacancy(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}
}

func TestSuccessDeleteVacancy(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("DELETE", "/api/vacancies/1", body)
	r = mux.SetURLVars(r, map[string]string{"vacancy_id": "1"})
	w := httptest.NewRecorder()

	h.DeleteVacancy(w, r)

	if w.Code != http.StatusNoContent {
		t.Error("Status code is not 204")
	}

	if len(h.vacancies) != 1 {
		t.Error("Vacancy is not deleted")
	}
}

func TestDeleteNullVacancy(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyVacancyHandler()

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("DELETE", "/api/vacancies/3", body)
	r = mux.SetURLVars(r, map[string]string{"vacancy_id": "3"})
	w := httptest.NewRecorder()

	h.DeleteVacancy(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status code is not 404")
	}
}