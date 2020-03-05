package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"joblessness/haha/handlers"
	"joblessness/haha/utils/cors"
	"joblessness/haha/utils/database"
	"log"
	"net/http"
)

func echoFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/users/echo")

	cors.Cors.PrivateApi(&w, r)

	params := mux.Vars(r)
	message := params["message"]
	fmt.Fprintf(w, "Hello %s!", message)
}


func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal haha error",
				})

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func StartRouter() {
	database.InitDatabase("username", "9730", "username")
	if err := database.OpenDatabase(); err != nil {
		log.Println(err.Error())
		return
	}
	defer database.CloseDatabase()

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	authApi := handlers.NewAuthHandler()
	vacancyApi := handlers.NewVacancyHandler()
	summaryApi := handlers.NewSummaryHandler()

	router.Use(RecoveryMiddleware)
	router.Use(cors.Cors.CorsMiddleware)
	router.HandleFunc("/echo/{message}", echoFunc)
	router.Methods("OPTIONS").HandlerFunc(cors.Cors.Preflight)

	// users

	router.HandleFunc("/users/login", authApi.Login).Methods("POST")
	router.HandleFunc("/users/check", authApi.Check).Methods("POST")
	router.HandleFunc("/users/logout", authApi.Logout).Methods("POST")
	router.HandleFunc("/users", authApi.Register).Methods("POST")

	router.HandleFunc("/user/{user_id}", authApi.GetUserPage).Methods("GET")
	router.HandleFunc("/users/{user_id}/avatar", authApi.SetAvatar).Methods("POST")
	router.HandleFunc("/user/{user_id}", authApi.ChangeUserInfo).Methods("POST")

	// vacancies

	router.HandleFunc("/vacancies", vacancyApi.CreateVacancy).Methods("POST")
	router.HandleFunc("/vacancies", vacancyApi.GetVacancies).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", vacancyApi.GetVacancy).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", vacancyApi.ChangeVacancy).Methods("PUT")
	router.HandleFunc("/vacancies/{vacancy_id}", vacancyApi.DeleteVacancy).Methods("DELETE")

	// summaries

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