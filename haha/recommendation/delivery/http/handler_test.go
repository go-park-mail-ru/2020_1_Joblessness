package recommendationHttp

//go:generate mockgen -destination=../../usecase/mock/recommendation.go -package=mock joblessness/haha/recommendation/interfaces UseCase

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/auth/usecase/mock"
	"joblessness/haha/middleware"
	"joblessness/haha/models/base"
	recommendationInterfaces "joblessness/haha/recommendation/interfaces"
	recommendationUseCaseMock "joblessness/haha/recommendation/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	router           *mux.Router
	mainMiddleware   *middleware.RecoveryHandler
	authMiddleware   *middleware.SessionHandler
	controller       *gomock.Controller
	authUseCase      *mock.MockAuthUseCase
	uc               *recommendationUseCaseMock.MockUseCase
	person           baseModels.Person
	personByte       *bytes.Buffer
	organization     baseModels.Organization
	organizationByte *bytes.Buffer
	cookie *http.Cookie
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/haha").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.router.Use(suite.mainMiddleware.LogMiddleware)

	suite.controller = gomock.NewController(suite.T())
	suite.uc = recommendationUseCaseMock.NewMockUseCase(suite.controller)
	suite.authUseCase = mock.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.authUseCase)

	suite.cookie = &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	RegisterHTTPEndpoints(suite.router, suite.authMiddleware, suite.uc)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestGetRecommendedVacancies() {
	suite.uc.EXPECT().
		GetRecommendedVacancies(uint64(1)).
		Return([]baseModels.Vacancy{}, nil).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/haha/recommendation", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetRecommendedVacanciesFailed() {
	suite.uc.EXPECT().
		GetRecommendedVacancies(uint64(1)).
		Return(nil, errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/haha/recommendation", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetRecommendedVacanciesNoRec() {
	suite.uc.EXPECT().
		GetRecommendedVacancies(uint64(1)).
		Return(nil, recommendationInterfaces.ErrNoRecommendation).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/haha/recommendation", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}