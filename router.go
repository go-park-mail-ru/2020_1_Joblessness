package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func echoFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users/logout")

	Cors.EnableCors(&w, r)

	params := mux.Vars(r)
	message := params["message"]
	fmt.Fprintf(w, "Hello %s!", message)
}

func StartRouter() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()//.StrictSlash(true)

	router.HandleFunc("/echo/{message}", echoFunc)
	//Марат
	//apiRouter.HandleFunc("/user/session", HandleRegistration).Methods("Put")
	//apiRouter.HandleFunc("/user/session", HandleAuthorisation).Methods("Post")

	authApi := NewAuthHandler()
	router.HandleFunc("/users/login", authApi.Login).Methods("POST")
	router.HandleFunc("/users/logout", authApi.Logout).Methods("POST")
	router.HandleFunc("/users", authApi.Register).Methods("POST")

	router.HandleFunc("/user/{id}", authApi.GetUserPage).Methods("GET")
	router.HandleFunc("/users/{id}/avatar", authApi.SetAvatar).Methods("PUT")

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