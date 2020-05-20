package authHttp

import (
	"github.com/gorilla/mux"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/middleware"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, uc authInterfaces.AuthUseCase) {
	h := NewHandler(uc)
	authRouter := router.PathPrefix("/users").Subrouter()
	organizationRouter := router.PathPrefix("/organizations").Subrouter()

	authRouter.HandleFunc("/login", h.Login).Methods("POST")
	authRouter.HandleFunc("/check", m.UserRequired(h.Check)).Methods("POST")
	authRouter.HandleFunc("/logout", h.Logout).Methods("POST")
	authRouter.HandleFunc("", h.RegisterPerson).Methods("POST")
	organizationRouter.HandleFunc("", h.RegisterOrg).Methods("POST")
}
