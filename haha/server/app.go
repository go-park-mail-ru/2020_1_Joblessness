package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"joblessness/haha/auth/delivery/http"
	"joblessness/haha/auth/interfaces"
	postgresAuth "joblessness/haha/auth/repository/postgres"
	usecaseAuth "joblessness/haha/auth/usecase"
	"joblessness/haha/middleware"
	httpSearch "joblessness/haha/search/delivery/http"
	searchInterfaces "joblessness/haha/search/interfaces"
	postgresSearch "joblessness/haha/search/repository/postgres"
	usecaseSearch "joblessness/haha/search/usecase"
	"joblessness/haha/summary/delivery/http"
	summaryInterfaces "joblessness/haha/summary/interfaces"
	postgresSummary "joblessness/haha/summary/repository/postgres"
	usecaseSummary "joblessness/haha/summary/usecase"
	"joblessness/haha/utils/cors"
	"joblessness/haha/utils/database"
	"joblessness/haha/vacancy/delivery/http"
	vacancyInterfaces "joblessness/haha/vacancy/interfaces"
	postgresVacancy "joblessness/haha/vacancy/repository/postgres"
	usecaseVacancy "joblessness/haha/vacancy/usecase"
	"net/http"
)

type App struct {
	httpServer *http.Server
	authUse    authInterfaces.AuthUseCase
	vacancyUse vacancyInterfaces.VacancyUseCase
	summaryUse summaryInterfaces.SummaryUseCase
	searchUse  searchInterfaces.SearchUseCase
	corsHandler *cors.CorsHandler
}

func NewApp(c *cors.CorsHandler) *App {
	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return nil
	}

	userRepo := postgresAuth.NewUserRepository(db)
	vacancyRepo := postgresVacancy.NewVacancyRepository(db)
	summaryRepo := postgresSummary.NewSummaryRepository(db)
	searchRepo := postgresSearch.NewAuthRepository(db)

	return &App{
		authUse: usecaseAuth.NewAuthUseCase(userRepo),
		vacancyUse: usecaseVacancy.NewVacancyUseCase(vacancyRepo),
		summaryUse: usecaseSummary.NewSummaryUseCase(summaryRepo),
		searchUse: usecaseSearch.NewSearchUseCase(searchRepo),
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

	// search
	httpSearch.RegisterHTTPEndpoints(router, app.searchUse)

	http.Handle("/", router)
	fmt.Println("Server started")
	http.ListenAndServe(":8001", router)
}