package server

import (
	"../cors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct{
	corsHandler *cors.CorsHandler
}

func NewServer() *Server {
	return &Server {
		corsHandler: cors.NewCorsHandler(),
	}
}

func (server *Server) AddOrigin(origin string) {
	server.corsHandler.AddOrigin(origin)
}

func (server *Server) echoFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/users/echo")

	server.corsHandler.PrivateApi(&w, r)

	params := mux.Vars(r)
	message := params["message"]
	fmt.Fprintf(w, "Hello %s!", message)
}

func (server *Server) Preflight(w http.ResponseWriter, req *http.Request) {
	server.corsHandler.PrivateApi(&w, req)
}

func (server *Server) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (server *Server) StartRouter() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()//.StrictSlash(true)

	router.Use(server.contentTypeMiddleware)
	router.HandleFunc("/echo/{message}", server.echoFunc)
	router.Methods("OPTIONS").HandlerFunc(server.Preflight)

	// users
	authApi := NewAuthHandler()

	router.HandleFunc("/users/login", authApi.Login).Methods("POST")
	router.HandleFunc("/users/check", authApi.Check).Methods("POST")
	router.HandleFunc("/users/logout", authApi.Logout).Methods("POST")
	router.HandleFunc("/users", authApi.Register).Methods("POST")

	router.HandleFunc("/user/{user_id}", authApi.GetUserPage).Methods("GET")
	router.HandleFunc("/users/{user_id}/avatar", authApi.SetAvatar).Methods("POST")
	router.HandleFunc("/user/{user_id}", authApi.ChangeUserInfo).Methods("POST")

	// vacancies
	vacancyApi := NewVacancyHandler()

	router.HandleFunc("/vacancies", vacancyApi.CreateVacancy).Methods("POST")
	router.HandleFunc("/vacancies", vacancyApi.GetVacancies).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", vacancyApi.GetVacancy).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", vacancyApi.ChangeVacancy).Methods("PUT")
	router.HandleFunc("/vacancies/{vacancy_id}", vacancyApi.DeleteVacancy).Methods("DELETE")

	// summaries
	summaryApi := NewSummaryHandler()

	router.HandleFunc("/summaries", summaryApi.CreateSummary).Methods("POST")
	router.HandleFunc("/summaries", summaryApi.GetSummaries).Methods("GET")
	router.HandleFunc("/summaries/{summary_id}", summaryApi.GetSummary).Methods("GET")
	router.HandleFunc("/summaries/{summary_id}", summaryApi.ChangeSummary).Methods("PUT")
	router.HandleFunc("/summaries/{summary_id}", summaryApi.DeleteSummary).Methods("DELETE")
	router.HandleFunc("/summaries/{summary_id}/print", summaryApi.PrintSummary).Methods("GET")
	router.HandleFunc("/user/{user_id}/summaries", summaryApi.GetUserSummaries).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Server started")
	http.ListenAndServe(":8001", router)
}