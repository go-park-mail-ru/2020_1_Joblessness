package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func StartRouter() {
	router := mux.NewRouter()//.StrictSlash(true)
	//apiRouter := router.PathPrefix("/api").Subrouter()

	//Марат
	//apiRouter.HandleFunc("/user/session", HandleRegistration).Methods("Put")
	//apiRouter.HandleFunc("/user/session", HandleAuthorisation).Methods("Post")

	authApi := NewAuthHandler()
	router.HandleFunc("/login", authApi.Login).Methods("POST")

	//Миша
	//apiRouter.HandleFunc("/user/{id}", HandleSetPrivateInfo).Methods("Put")
	//apiRouter.HandleFunc("/user/{id}", HandleGetPrivateInfo).Methods("Get")
	//apiRouter.HandleFunc("/avatar/{id}", HandleSetAvatar).Methods("Put")
	//apiRouter.HandleFunc("/avatar/{id}", HandleGetAvatar).Methods("Get")

	//Huvalk
	//apiRouter.HandleFunc("/resume", HadleCreateResume).Methods("Post")
	//apiRouter.HandleFunc("/resume/{id}", HandleChangeResume).Methods("Put")
	//apiRouter.HandleFunc("/resume/{id}", HandleGetResume).Methods("Get")
	//apiRouter.HandleFunc("/resume/{id}", HandleRemoveResume).Methods("Delete")

	//Huvalk
	//apiRouter.HandleFunc("/vacancy", HandleCreateVacancy).Methods("Post")
	//apiRouter.HandleFunc("/vacancy/{id}", HandleChangeVacancy).Methods("Put")
	//apiRouter.HandleFunc("/vacancy/{id}", HandleGetVacancy).Methods("Get")
	//apiRouter.HandleFunc("/vacancy/{id}", HandleRemoveVacancy).Methods("Delete")

	//Сережа М
	//apiRouter.HandleFunc("/resume/short/{id}", HandleGetShortResume).Methods("Get")
	//apiRouter.HandleFunc("/vacancy/short/{id}", HandleGetShortVacancy).Methods("Get")
	//apiRouter.HandleFunc("/vacancys/{from}/{to}", HandleGetShortVacancy).Methods("Get")

	http.Handle("/", router)
	fmt.Println("Server started")
	http.ListenAndServe(":8000", router)
}