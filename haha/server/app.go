package server

import (
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/microcosm-cc/bluemonday"
	"joblessness/haha/auth/delivery/http"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/auth/repository/postgres"
	"joblessness/haha/auth/usecase"
	"joblessness/haha/interview/delivery/http"
	"joblessness/haha/interview/interfaces"
	"joblessness/haha/interview/repository/postgres"
	"joblessness/haha/interview/usecase"
	"joblessness/haha/middleware"
	"joblessness/haha/search/delivery/http"
	"joblessness/haha/search/interfaces"
	"joblessness/haha/search/repository/postgres"
	"joblessness/haha/search/usecase"
	"joblessness/haha/summary/delivery/http"
	"joblessness/haha/summary/interfaces"
	"joblessness/haha/summary/repository/postgres"
	"joblessness/haha/summary/usecase"
	"joblessness/haha/user/delivery/http"
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
	httpServer  *http.Server
	userUse     userInterfaces.UserUseCase
	authUse     authInterfaces.AuthUseCase
	vacancyUse  vacancyInterfaces.VacancyUseCase
	summaryUse  summaryInterfaces.SummaryUseCase
	searchUse   searchInterfaces.SearchUseCase
	interviewUse   interviewInterfaces.InterviewUseCase
	corsHandler *middleware.CorsHandler
}

func NewApp(c *middleware.CorsHandler) *App {
	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return nil
	}

	userRepo := userPostgres.NewUserRepository(db)
	authRepo := authPostgres.NewAuthRepository(db)
	vacancyRepo := vacancyPostgres.NewVacancyRepository(db)
	summaryRepo := summaryPostgres.NewSummaryRepository(db)
	searchRepo := searchPostgres.NewSearchRepository(db)
	interviewRepo := interviewPostgres.NewInterviewRepository(db)
	policy := bluemonday.UGCPolicy()

	return &App{
		userUse:     userUseCase.NewUserUseCase(userRepo, policy),
		authUse:     authUseCase.NewAuthUseCase(authRepo),
		vacancyUse:  vacancyUseCase.NewVacancyUseCase(vacancyRepo, policy),
		summaryUse:  summaryUseCase.NewSummaryUseCase(summaryRepo, policy),
		searchUse:   searchUseCase.NewSearchUseCase(searchRepo, policy),
		interviewUse:   interviewUseCase.NewInterviewUseCase(interviewRepo, policy),
		corsHandler: c,
	}
}

func (app *App) StartRouter() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	m := middleware.NewMiddleware()
	mAuth := middleware.NewAuthMiddleware(app.authUse)

	router.Use(m.RecoveryMiddleware)
	//router.Use(app.corsHandler.CorsMiddleware)
	router.Use(m.LogMiddleware)
	router.Methods("OPTIONS").HandlerFunc(app.corsHandler.Preflight)

	// auth
	authHttp.RegisterHTTPEndpoints(router, mAuth, app.authUse)

	// users
	userHttp.RegisterHTTPEndpoints(router, mAuth, app.userUse)

	// vacancies
	vacancyHttp.RegisterHTTPEndpoints(router, mAuth, app.vacancyUse)

	// summaries
	summaryHttp.RegisterHTTPEndpoints(router, mAuth, app.summaryUse)

	// search
	searchHttp.RegisterHTTPEndpoints(router, app.searchUse)

	// interview
	interviewHttp.RegisterHTTPEndpoints(router, mAuth, app.interviewUse)

	http.Handle("/", router)
	golog.Info("Server started")
	err := http.ListenAndServe(":8001", router)
	if err != nil {
		golog.Error("Server failed")
	}
}
