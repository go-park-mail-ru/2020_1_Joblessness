package httpAuth

import (
	"github.com/gorilla/mux"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/middleware"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, uc authInterfaces.AuthUseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/users/login", h.Login).Methods("POST")
	router.HandleFunc("/users/check", m.UserRequired(h.Check)).Methods("POST")
	router.HandleFunc("/users/logout", h.Logout).Methods("POST")
	router.HandleFunc("/users", h.RegisterPerson).Methods("POST")
	router.HandleFunc("/organizations", h.RegisterOrg).Methods("POST")
	router.HandleFunc("/users/{user_id:[0-9]+}", h.GetPerson).Methods("GET")
	router.HandleFunc("/users/{user_id:[0-9]+}", m.UserRequired(h.ChangePerson)).Methods("PUT")
	router.HandleFunc("/organizations/{user_id:[0-9]+}", h.GetOrganization).Methods("GET")
	router.HandleFunc("/organizations/{user_id:[0-9]+}", m.UserRequired(h.ChangeOrganization)).Methods("PUT")
	router.HandleFunc("/organizations", h.GetListOfOrgs).Methods("GET")
	router.HandleFunc("/users/{user_id:[0-9]+}/avatar", m.UserRequired(h.SetAvatar)).Methods("POST")
	router.HandleFunc("/users/{user_id:[0-9]+}/like", m.UserRequired(h.LikeUser)).Methods("POST")
	router.HandleFunc("/users/{user_id:[0-9]+}/like", m.UserRequired(h.LikeExists)).Methods("GET")
	router.HandleFunc("/users/{user_id:[0-9]+}/favorite", m.UserRequired(h.GetUserFavorite)).Methods("GET")
}