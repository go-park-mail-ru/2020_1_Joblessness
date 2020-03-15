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
	router.HandleFunc("/users", h.RegisterPerson).Methods("POST")
	router.HandleFunc("/organizations", h.RegisterOrg).Methods("POST")
	router.HandleFunc("/users/{user_id:[0-9]+}", h.GetPerson).Methods("GET")
	router.HandleFunc("/users/{user_id:[0-9]+}", h.ChangePerson).Methods("PUT")
	router.HandleFunc("/organizations/{user_id:[0-9]+}", h.GetOrganization).Methods("GET")
	router.HandleFunc("/organizations/{user_id:[0-9]+}", h.ChangeOrganization).Methods("PUT")
}