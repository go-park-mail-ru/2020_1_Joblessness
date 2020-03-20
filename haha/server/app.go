package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"joblessness/haha/auth"
	"joblessness/haha/auth/delivery/http"
	postgresAuth "joblessness/haha/auth/repository/postgres"
	usecaseAuth "joblessness/haha/auth/usecase"
	"joblessness/haha/middleware"
	"joblessness/haha/summary"
	"joblessness/haha/summary/delivery/http"
	postgresSummary "joblessness/haha/summary/repository/postgres"
	usecaseSummary "joblessness/haha/summary/usecase"
	"joblessness/haha/utils/cors"
	"joblessness/haha/utils/database"
	"joblessness/haha/utils/seed"
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
	authUse auth.AuthUseCase
	vacancyUse vacancy.VacancyUseCase
	summaryUse summary.SummaryUseCase
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
	summaryRepo := postgresSummary.NewSummaryRepository(db)

	return &App{
		authUse: usecaseAuth.NewAuthUseCase(userRepo),
		vacancyUse: usecaseVacancy.NewVacancyUseCase(vacancyRepo),
		summaryUse: usecaseSummary.NewSummaryUseCase(summaryRepo),
		corsHandler: c,
	}
}

func (app *App) StartRouter() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	m := middleware.NewMiddleware()
	mAuth := middleware.NewAuthMiddleware(app.authUse)

	router.Use(m.RecoveryMiddleware)
	router.Use(app.corsHandler.CorsMiddleware)
	router.Use(m.LogMiddleware)
	router.Methods("OPTIONS").HandlerFunc(app.corsHandler.Preflight)

	// users
	httpAuth.RegisterHTTPEndpoints(router, mAuth, app.authUse)

	// vacancies
	httpVacancy.RegisterHTTPEndpoints(router, mAuth, app.vacancyUse)

	// summaries
	httpSummary.RegisterHTTPEndpoints(router, mAuth, app.summaryUse)

	seeder := seed.NewSeeder(&app.authUse)
	err := seeder.Seed()
	if err != nil {
		golog.Error(err.Error())
	}

	http.Handle("/", router)
	fmt.Println("Server started")
	http.ListenAndServe(":8001", router)
}