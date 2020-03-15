package httpVacancy

import (
	"github.com/gorilla/mux"
	"joblessness/haha/vacancy"
)

func RegisterHTTPEndpoints(router *mux.Router, uc vacancy.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/vacancies", h.CreateVacancy).Methods("POST")
	router.HandleFunc("/vacancies", h.GetVacancies).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", h.GetVacancy).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", h.ChangeVacancy).Methods("PUT")
	router.HandleFunc("/vacancies/{vacancy_id}", h.DeleteVacancy).Methods("DELETE")
}
