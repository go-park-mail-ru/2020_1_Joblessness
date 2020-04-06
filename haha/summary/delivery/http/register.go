package httpSummary

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/summary/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, uc summaryInterfaces.SummaryUseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/summaries", m.PersonRequired(h.CreateSummary)).Methods("POST")
	router.HandleFunc("/summaries", h.GetSummaries).Methods("GET")
	router.HandleFunc("/summaries/{summary_id:[0-9]+}", h.GetSummary).Methods("GET")
	router.HandleFunc("/summaries/{summary_id:[0-9]+}", m.PersonRequired(h.ChangeSummary)).Methods("PUT")
	router.HandleFunc("/summaries/{summary_id:[0-9]+}", m.PersonRequired(h.DeleteSummary)).Methods("DELETE")
	router.HandleFunc("/summaries/{summary_id:[0-9]+}/print", h.PrintSummary).Methods("GET")
	router.HandleFunc("/users/{user_id:[0-9]+}/summaries", h.GetUserSummaries).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id:[0-9]+}/response", m.PersonRequired(h.SendSummary)).Methods("POST")
	router.HandleFunc("/summaries/{summary_id:[0-9]+}/response", m.OrganizationRequired(h.ResponseSummary)).Methods("PUT")
	router.HandleFunc("/organizations/{user_id:[0-9]+}/response", h.GetOrgSendSummaries).Methods("GET")
	router.HandleFunc("/users/{user_id:[0-9]+}/response", h.GetUserSendSummaries).Methods("GET")
}
