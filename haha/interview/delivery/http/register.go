package interviewHttp

import (
	"github.com/gorilla/mux"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	"joblessness/haha/middleware"
	"joblessness/haha/utils/chat"
)

func RegisterHTTPEndpoints(router *mux.Router,
	m *middleware.SessionHandler,
	uc interviewInterfaces.InterviewUseCase,
	room chat.Room) {
	h := NewHandler(uc, room)
	chatRouter := router.PathPrefix("/chat").Subrouter()

	router.HandleFunc("/summaries/{summary_id:[0-9]+}/response", m.OrganizationRequired(h.ResponseSummary)).Methods("PUT")
	chatRouter.HandleFunc("", m.UserRequired(h.EnterChat)).Methods("PUT")
	chatRouter.HandleFunc("/conversation/{user_id:[0-9]+}", m.UserRequired(h.History)).Methods("GET")
	chatRouter.HandleFunc("/conversation", m.UserRequired(h.GetConversations)).Methods("GET")
}