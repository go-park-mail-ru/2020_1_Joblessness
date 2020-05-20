package userHttp

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/user/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, uc userInterfaces.UserUseCase) {
	h := NewHandler(uc)
	usersRouter := router.PathPrefix("/users").Subrouter()
	organizationRouter := router.PathPrefix("/organizations").Subrouter()

	usersRouter.HandleFunc("/{user_id:[0-9]+}", h.GetPerson).Methods("GET")
	usersRouter.HandleFunc("/{user_id:[0-9]+}", m.UserRequired(h.ChangePerson)).Methods("PUT")
	organizationRouter.HandleFunc("/{user_id:[0-9]+}", h.GetOrganization).Methods("GET")
	organizationRouter.HandleFunc("/{user_id:[0-9]+}", m.UserRequired(h.ChangeOrganization)).Methods("PUT")
	organizationRouter.HandleFunc("", h.GetListOfOrgs).Methods("GET")
	usersRouter.HandleFunc("/{user_id:[0-9]+}/avatar", m.UserRequired(h.SetAvatar)).Methods("POST")
	usersRouter.HandleFunc("/{user_id:[0-9]+}/like", m.UserRequired(h.LikeUser)).Methods("POST")
	usersRouter.HandleFunc("/{user_id:[0-9]+}/like", m.UserRequired(h.LikeExists)).Methods("GET")
	usersRouter.HandleFunc("/{user_id:[0-9]+}/favorite", m.UserRequired(h.GetUserFavorite)).Methods("GET")
}
