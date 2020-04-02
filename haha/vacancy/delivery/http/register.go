package httpVacancy

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/vacancy/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, uc vacancyInterfaces.VacancyUseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/vacancies", m.OrganizationRequired(h.CreateVacancy)).Methods("POST")
	router.HandleFunc("/vacancies", h.GetVacancies).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", h.GetVacancy).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", m.OrganizationRequired(h.ChangeVacancy)).Methods("PUT")
	router.HandleFunc("/vacancies/{vacancy_id}", m.OrganizationRequired(h.DeleteVacancy) ).Methods("DELETE")
	router.HandleFunc("/organizations/{organization_id}/vacancies", h.GetVacancies ).Methods("GET")
}
