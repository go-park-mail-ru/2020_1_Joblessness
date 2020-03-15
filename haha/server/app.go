package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"joblessness/haha/auth"
	"joblessness/haha/auth/delivery/http"
	postgresAuth "joblessness/haha/auth/repository/postgres"
	usecaseAuth "joblessness/haha/auth/usecase"
	"joblessness/haha/handlers"
	"joblessness/haha/middleware"
	"joblessness/haha/utils/cors"
	"joblessness/haha/utils/database"
	"joblessness/haha/vacancy"
	"joblessness/haha/vacancy/delivery/http"
	postgresVacancy "joblessness/haha/vacancy/repository/postgres"
	usecaseVacancy "joblessness/haha/vacancy/usecase"
	"log"
	"net/http"
	"os"
)

type App struct {
	httpServer *http.Server
	authUse auth.UseCase
	vacancyUse vacancy.UseCase
	corsHandler *cors.CorsHandler
}

func NewApp(c *cors.CorsHandler) *App {
	database.InitDatabase(os.Getenv("HAHA_DB_USER"), os.Getenv("HAHA_DB_PASSWORD"), os.Getenv("HAHA_DB_NAME"))
	if err := database.OpenDatabase(); err != nil {
		log.Println(err.Error())
		return nil
	}
	db := database.GetDatabase()


	userRepo := postgresAuth.NewUserRepository(db)
	vacancyRepo := postgresVacancy.NewVacancyRepository(db)

	return &App{
		authUse: usecaseAuth.NewAuthUseCase(userRepo),
		vacancyUse: usecaseVacancy.NewVacancyUseCase(vacancyRepo),
		corsHandler: c,
	}
}

func (app *App) StartRouter() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	authApi := handlers.NewAuthHandler()
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
	httpVacancy.RegisterHTTPEndpoints(router, app.vacancyUse)

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