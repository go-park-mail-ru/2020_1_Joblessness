package server

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc"
	"joblessness/haha/auth/delivery/http"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/auth/repository/grpc"
	"joblessness/haha/auth/usecase"
	interviewHttp "joblessness/haha/interview/delivery/http"
	"joblessness/haha/interview/interfaces"
	"joblessness/haha/interview/repository/grpc"
	"joblessness/haha/interview/usecase"
	"joblessness/haha/middleware"
	"joblessness/haha/recommendation/delivery/http"
	"joblessness/haha/recommendation/interfaces"
	"joblessness/haha/recommendation/repository/postgres"
	"joblessness/haha/recommendation/usecase"
	"joblessness/haha/search/delivery/http"
	"joblessness/haha/search/interfaces"
	"joblessness/haha/search/repository/grpc"
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
	httpServer        *http.Server
	userUse           userInterfaces.UserUseCase
	authUse        		authInterfaces.AuthUseCase
	vacancyUse        vacancyInterfaces.VacancyUseCase
	summaryUse        summaryInterfaces.SummaryUseCase
	searchUse         searchInterfaces.SearchUseCase
	recommendationUse recommendationInterfaces.UseCase
	interviewUse      interviewInterfaces.InterviewUseCase
	corsHandler       *middleware.CorsHandler
}

func NewApp(c *middleware.CorsHandler) *App {
	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return nil
	}
	err = db.Ping()
	if err != nil {
		golog.Error("DB: ", err.Error())
		return nil
	}

	searchConn, err := grpc.Dial(
		"127.0.0.1:8002",
		grpc.WithInsecure(),
	)
	if err != nil {
		golog.Fatal("cant connect to grpc")
	}

	interviewConn, err := grpc.Dial(
		"127.0.0.1:8003",
		grpc.WithInsecure(),
	)
	if err != nil {
		golog.Fatal("cant connect to grpc")
	}

	authConn, err := grpc.Dial(
		"127.0.0.1:8003",
		grpc.WithInsecure(),
	)
	if err != nil {
		golog.Fatal("can't connect to auth")
	}

	userRepo := userPostgres.NewUserRepository(db)
	authRepo := authGrpcRepository.NewRepository(authConn)
	vacancyRepo := vacancyPostgres.NewVacancyRepository(db)
	summaryRepo := summaryPostgres.NewSummaryRepository(db)
	searchRepo := searchGrpc.NewSearchGrpcRepository(searchConn)
	recommendationRepo := recommendationPostgres.NewRepository(db, vacancyRepo)
	interviewRepo := interviewGrpc.NewInterviewGrpcRepository(interviewConn)
	policy := bluemonday.UGCPolicy()

	interviewUse, room := interviewUseCase.NewInterviewUseCase(interviewRepo, policy)

	return &App{
		userUse:           userUseCase.NewUserUseCase(userRepo, policy),
		authUse:           authUseCase.NewAuthUseCase(authRepo),
		vacancyUse:        vacancyUseCase.NewVacancyUseCase(vacancyRepo, room, policy),
		summaryUse:        summaryUseCase.NewSummaryUseCase(summaryRepo, policy),
		searchUse:         searchUseCase.NewSearchUseCase(searchRepo, policy),
		recommendationUse: recommendationUseCase.NewUseCase(recommendationRepo),
		interviewUse:      interviewUse,
		corsHandler:       c,
	}
}

var (
	noCors = flag.Bool("no-cors", false, "disable cors")
	port   = flag.Uint64("p", 8001, "port")
)

func (app *App) StartRouter() {
	flag.Parse()

	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	m := middleware.NewMiddleware()
	mAuth := middleware.NewAuthMiddleware(app.authUse)

	router.Use(m.RecoveryMiddleware)
	if !*noCors {
		router.Use(app.corsHandler.CorsMiddleware)
	}
	router.Use(m.LogMiddleware)
	router.Methods("OPTIONS").HandlerFunc(app.corsHandler.Preflight)

	authHttp.RegisterHTTPEndpoints(router, mAuth, app.authUse)
	userHttp.RegisterHTTPEndpoints(router, mAuth, app.userUse)
	vacancyHttp.RegisterHTTPEndpoints(router, mAuth, app.vacancyUse)
	summaryHttp.RegisterHTTPEndpoints(router, mAuth, app.summaryUse)
	searchHttp.RegisterHTTPEndpoints(router, app.searchUse)
	recommendationHttp.RegisterHTTPEndpoints(router, mAuth, app.recommendationUse)
	interviewHttp.RegisterHTTPEndpoints(router, mAuth, app.interviewUse)

	http.Handle("/", router)
	golog.Infof("Server started at port :%d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router)
	if err != nil {
		golog.Error("Server failed")
	}
}
