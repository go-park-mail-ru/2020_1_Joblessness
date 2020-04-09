package server

import (
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/microcosm-cc/bluemonday"
	httpAuth "joblessness/haha/auth/delivery/http"
	"joblessness/haha/auth/interfaces"
	postgresAuth "joblessness/haha/auth/repository/postgres"
	authUseCase "joblessness/haha/auth/usecase"
	"joblessness/haha/middleware"
	"joblessness/haha/search/delivery/http"
	"joblessness/haha/search/interfaces"
	"joblessness/haha/search/repository/postgres"
	"joblessness/haha/search/usecase"
	"joblessness/haha/summary/delivery/http"
	"joblessness/haha/summary/interfaces"
	"joblessness/haha/summary/repository/postgres"
	"joblessness/haha/summary/usecase"
	httpUser "joblessness/haha/user/delivery/http"
	"joblessness/haha/user/interfaces"
	"joblessness/haha/user/repository/postgres"
	"joblessness/haha/user/usecase"
	"joblessness/haha/utils/database"
	"joblessness/haha/vacancy/delivery/http"
	"joblessness/haha/vacancy/interfaces"
	"joblessness/haha/vacancy/repository/postgres"
	"joblessness/haha/vacancy/usecase"
	"net/http"
)

type App struct {
	httpServer *http.Server
	userUse    userInterfaces.UserUseCase
	authUse authInterfaces.AuthUseCase
	vacancyUse vacancyInterfaces.VacancyUseCase
	summaryUse summaryInterfaces.SummaryUseCase
	searchUse  searchInterfaces.SearchUseCase
	corsHandler *middleware.CorsHandler
}

func NewApp(c *middleware.CorsHandler) *App {
	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return nil
	}

	userRepo := postgresUser.NewUserRepository(db)
	authRepo := postgresAuth.NewAuthRepository(db)
	vacancyRepo := vacancyRepoPostgres.NewVacancyRepository(db)
	summaryRepo := summaryRepoPostgres.NewSummaryRepository(db)
	searchRepo := searchRepoPostgres.NewSearchRepository(db)
	policy := bluemonday.UGCPolicy()

	return &App{
		userUse: userUseCase.NewUserUseCase(userRepo, policy),
		authUse: authUseCase.NewAuthUseCase(authRepo),
		vacancyUse: vacancyUseCase.NewVacancyUseCase(vacancyRepo, policy),
		summaryUse: summaryUseCase.NewSummaryUseCase(summaryRepo, policy),
		searchUse: searchUseCase.NewSearchUseCase(searchRepo, policy),
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

	// auth
	httpAuth.RegisterHTTPEndpoints(router, mAuth, app.authUse)

	// users
	httpUser.RegisterHTTPEndpoints(router, mAuth, app.userUse)

	// vacancies
	httpVacancy.RegisterHTTPEndpoints(router, mAuth, app.vacancyUse)

	// summaries
	httpSummary.RegisterHTTPEndpoints(router, mAuth, app.summaryUse)

	// search
	httpSearch.RegisterHTTPEndpoints(router, app.searchUse)

	http.Handle("/", router)
	golog.Info("Server started")
	err := http.ListenAndServe(":8001", router)
	if err != nil {
		golog.Error("Server failed")
	}
}