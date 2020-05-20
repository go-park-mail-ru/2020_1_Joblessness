package vacancyHttp

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/vacancy/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, uc vacancyInterfaces.VacancyUseCase) {
	h := NewHandler(uc)
	vacanciesRouter := router.PathPrefix("/vacancies").Subrouter()

	vacanciesRouter.HandleFunc("", m.OrganizationRequired(h.CreateVacancy)).Methods("POST")
	vacanciesRouter.HandleFunc("", h.GetVacancies).Methods("GET")
	vacanciesRouter.HandleFunc("/{vacancy_id:[0-9]+}", h.GetVacancy).Methods("GET")
	vacanciesRouter.HandleFunc("/{vacancy_id:[0-9]+}", m.OrganizationRequired(h.ChangeVacancy)).Methods("PUT")
	vacanciesRouter.HandleFunc("/{vacancy_id:[0-9]+}", m.OrganizationRequired(h.DeleteVacancy)).Methods("DELETE")
	router.HandleFunc("/organizations/{organization_id:[0-9]+}/vacancies", h.GetOrgVacancies).Methods("GET")
}
