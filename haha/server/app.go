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
	"joblessness/haha/prometheus"
	"joblessness/haha/recommend/delivery/http"
	"joblessness/haha/recommend/interfaces"
	"joblessness/haha/recommend/repository/postgres"
	"joblessness/haha/recommend/usecase"
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
	"joblessness/haha/utils/chat"
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
	authUse           authInterfaces.AuthUseCase
	vacancyUse        vacancyInterfaces.VacancyUseCase
	summaryUse        summaryInterfaces.SummaryUseCase
	searchUse         searchInterfaces.SearchUseCase
	recommendationUse recommendInterfaces.RecommendUseCase
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
		"127.0.0.1:8004",
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
	recommendationRepo := recommendPostgres.NewRecommendRepository(db, vacancyRepo)
	interviewRepo := interviewGrpc.NewInterviewGrpcRepository(interviewConn)
	policy := bluemonday.UGCPolicy()

	interviewUse := interviewUseCase.NewInterviewUseCase(interviewRepo, policy)
	room := chat.NewRoom(interviewUse)
	interviewUse.EnableRoom(room)

	return &App{
		userUse:           userUseCase.NewUserUseCase(userRepo, policy),
		authUse:           authUseCase.NewAuthUseCase(authRepo),
		vacancyUse:        vacancyUseCase.NewVacancyUseCase(vacancyRepo, room, policy),
		summaryUse:        summaryUseCase.NewSummaryUseCase(summaryRepo, policy),
		searchUse:         searchUseCase.NewSearchUseCase(searchRepo, policy),
		recommendationUse: recommendUseCase.NewUseCase(recommendationRepo),
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

	router := mux.NewRouter()

	commonRouter := router.PathPrefix("/api").Subrouter()
	wsRouter := router.PathPrefix("/api").Subrouter()

	m := middleware.NewMiddleware()
	mAuth := middleware.NewAuthMiddleware(app.authUse)

	router.Use(m.RecoveryMiddleware)
	if !*noCors {
		router.Use(app.corsHandler.CorsMiddleware)
		golog.Info("Cors enabled")
	}
	commonRouter.Use(m.LogMiddleware)
	router.Methods("OPTIONS").HandlerFunc(app.corsHandler.Preflight)

	authHttp.RegisterHTTPEndpoints(commonRouter, mAuth, app.authUse)
	userHttp.RegisterHTTPEndpoints(commonRouter, mAuth, app.userUse)
	vacancyHttp.RegisterHTTPEndpoints(commonRouter, mAuth, app.vacancyUse)
	summaryHttp.RegisterHTTPEndpoints(commonRouter, mAuth, app.summaryUse)
	searchHttp.RegisterHTTPEndpoints(commonRouter, app.searchUse)
	recommendHttp.RegisterHTTPEndpoints(commonRouter, mAuth, app.recommendationUse)
	interviewHttp.RegisterHTTPEndpoints(commonRouter, wsRouter, mAuth, app.interviewUse)

	prometheus.RegisterPrometheus(commonRouter)

	http.Handle("/", router)
	golog.Infof("Server started at port :%d", *port)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", *port),
		"/etc/letsencrypt/live/hahao.ru/fullchain.pem",
		"/etc/letsencrypt/live/hahao.ru/privkey.pem",
		nil)
	//err := http.ListenAndServe(":8001", nil)
	if err != nil {
		golog.Error("Server haha failed: ", err)
	}
}
