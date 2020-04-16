package summaryHttp

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/summary/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, uc summaryInterfaces.SummaryUseCase) {
	h := NewHandler(uc)
	summariesRouter := router.PathPrefix("/summaries").Subrouter()
	usersRouter := router.PathPrefix("/users").Subrouter()

	summariesRouter.HandleFunc("", m.PersonRequired(h.CreateSummary)).Methods("POST")
	summariesRouter.HandleFunc("", h.GetSummaries).Methods("GET")
	summariesRouter.HandleFunc("/{summary_id:[0-9]+}", h.GetSummary).Methods("GET")
	summariesRouter.HandleFunc("/{summary_id:[0-9]+}", m.PersonRequired(h.ChangeSummary)).Methods("PUT")
	summariesRouter.HandleFunc("/{summary_id:[0-9]+}", m.PersonRequired(h.DeleteSummary)).Methods("DELETE")
	summariesRouter.HandleFunc("/{summary_id:[0-9]+}/print", h.PrintSummary).Methods("GET")
	usersRouter.HandleFunc("/{user_id:[0-9]+}/summaries", h.GetUserSummaries).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id:[0-9]+}/response", m.PersonRequired(h.SendSummary)).Methods("POST")
	summariesRouter.HandleFunc("/{summary_id:[0-9]+}/response", m.OrganizationRequired(h.ResponseSummary)).Methods("PUT")
	router.HandleFunc("/organizations/{user_id:[0-9]+}/response", h.GetOrgSendSummaries).Methods("GET")
	usersRouter.HandleFunc("/{user_id:[0-9]+}/response", h.GetUserSendSummaries).Methods("GET")
	summariesRouter.HandleFunc("/{summary_id:[0-9]+}/mail", m.PersonRequired(h.SendSummaryByMail)).Methods("POST")
}
