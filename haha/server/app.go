package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"joblessness/haha/auth"
	"joblessness/haha/auth/delivery/http"
	"joblessness/haha/auth/repository/postgres"
	"joblessness/haha/auth/usecase"
	"joblessness/haha/handlers"
	"joblessness/haha/middleware"
	"joblessness/haha/utils/cors"
	"joblessness/haha/utils/database"
	"log"
	"net/http"
	"os"
)

type App struct {
	httpServer *http.Server
	authUse auth.UseCase
	corsHandler *cors.CorsHandler
}

func NewApp(c *cors.CorsHandler) *App {
	database.InitDatabase(os.Getenv("HAHA_DB_USER"), os.Getenv("HAHA_DB_PASSWORD"), os.Getenv("HAHA_DB_NAME"))
	if err := database.OpenDatabase(); err != nil {
		log.Println(err.Error())
		return nil
	}
	db := database.GetDatabase()


	userRepo := postgres.NewUserRepository(db)

	return &App{
		authUse: usecase.NewAuthUseCase(userRepo),
		corsHandler: c,
	}
}

func (app *App) StartRouter() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	authApi := handlers.NewAuthHandler()
	vacancyApi := handlers.NewVacancyHandler()
	summaryApi := handlers.NewSummaryHandler()

	m := middleware.NewMiddleware(app.authUse)

	router.Use(m.RecoveryMiddleware)
	router.Use(app.corsHandler.CorsMiddleware)
	router.Use(m.LogMiddleware)
	router.Methods("OPTIONS").HandlerFunc(app.corsHandler.Preflight)

	// users
	httpAuth.RegisterHTTPEndpoints(router, m, app.authUse)

	router.HandleFunc("/users/{user_id}/avatar", authApi.SetAvatar).Methods("POST")

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