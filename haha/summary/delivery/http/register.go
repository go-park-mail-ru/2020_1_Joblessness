package httpSummary

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/summary"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.AuthMiddleware, uc summary.SummaryUseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/summaries", m.CheckAuth(h.CreateSummary)).Methods("POST")
	router.HandleFunc("/summaries", h.GetSummaries).Methods("GET")
	router.HandleFunc("/summaries/{summary_id}", h.GetSummary).Methods("GET")
	router.HandleFunc("/summaries/{summary_id}", m.CheckAuth(h.ChangeSummary)).Methods("PUT")
	router.HandleFunc("/summaries/{summary_id}", m.CheckAuth(h.DeleteSummary)).Methods("DELETE")
	router.HandleFunc("/summaries/{summary_id}/print", h.PrintSummary).Methods("GET")
	router.HandleFunc("/user/{user_id}/summaries", h.GetUserSummaries).Methods("GET")
}
