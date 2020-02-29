package routers

import (
	_handlers "../../handlers"
	_cors "../cors"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func echoFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/users/echo")

	_cors.Cors.PrivateApi(&w, r)

	params := mux.Vars(r)
	message := params["message"]
	fmt.Fprintf(w, "Hello %s!", message)
}


func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				log.Println(err) // May be log this error? Send to sentry?

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func StartRouter() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	router.Use(RecoveryMiddleware)
	router.HandleFunc("/echo/{message}", echoFunc)
	router.Methods("OPTIONS").HandlerFunc(_cors.Cors.Preflight)

	// users
	authApi := _handlers.NewAuthHandler()

	router.HandleFunc("/users/login", _cors.Cors.CorsMiddleware(authApi.Login)).Methods("POST")
	router.HandleFunc("/users/check", _cors.Cors.CorsMiddleware(authApi.Check)).Methods("POST")
	router.HandleFunc("/users/logout", _cors.Cors.CorsMiddleware(authApi.Logout)).Methods("POST")
	router.HandleFunc("/users", _cors.Cors.CorsMiddleware(authApi.Register)).Methods("POST")

	router.HandleFunc("/user/{user_id}", _cors.Cors.CorsMiddleware(authApi.GetUserPage)).Methods("GET")
	router.HandleFunc("/users/{user_id}/avatar", _cors.Cors.CorsMiddleware(authApi.SetAvatar)).Methods("POST")
	router.HandleFunc("/user/{user_id}", _cors.Cors.CorsMiddleware(authApi.ChangeUserInfo)).Methods("POST")

	// vacancies
	vacancyApi := _handlers.NewVacancyHandler()

	router.HandleFunc("/vacancies", _cors.Cors.CorsMiddleware(vacancyApi.CreateVacancy)).Methods("POST")
	router.HandleFunc("/vacancies", _cors.Cors.CorsMiddleware(vacancyApi.GetVacancies)).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", _cors.Cors.CorsMiddleware(vacancyApi.GetVacancy)).Methods("GET")
	router.HandleFunc("/vacancies/{vacancy_id}", _cors.Cors.CorsMiddleware(vacancyApi.ChangeVacancy)).Methods("PUT")
	router.HandleFunc("/vacancies/{vacancy_id}", _cors.Cors.CorsMiddleware(vacancyApi.DeleteVacancy)).Methods("DELETE")

	// summaries
	summaryApi := _handlers.NewSummaryHandler()

	router.HandleFunc("/summaries", _cors.Cors.CorsMiddleware(summaryApi.CreateSummary)).Methods("POST")
	router.HandleFunc("/summaries", _cors.Cors.CorsMiddleware(summaryApi.GetSummaries)).Methods("GET")
	router.HandleFunc("/summaries/{summary_id}", _cors.Cors.CorsMiddleware(summaryApi.GetSummary)).Methods("GET")
	router.HandleFunc("/summaries/{summary_id}", _cors.Cors.CorsMiddleware(summaryApi.ChangeSummary)).Methods("PUT")
	router.HandleFunc("/summaries/{summary_id}", _cors.Cors.CorsMiddleware(summaryApi.DeleteSummary)).Methods("DELETE")
	router.HandleFunc("/user/{user_id}/summaries", _cors.Cors.CorsMiddleware(summaryApi.GetUserSummaries)).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Server started")
	http.ListenAndServe(":8001", router)
}