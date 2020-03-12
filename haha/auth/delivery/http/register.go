package httpAuth

import (
	"github.com/gorilla/mux"
	"joblessness/haha/auth"
)

func RegisterHTTPEndpoints(router *mux.Router, uc auth.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/users/login", h.Login).Methods("POST")
	router.HandleFunc("/users/check", h.Check).Methods("POST")
	router.HandleFunc("/users/logout", h.Logout).Methods("POST")
	router.HandleFunc("/users", h.Register).Methods("POST")
}