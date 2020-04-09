package httpSearch

//go:generate mockgen -destination=../../usecase/mock/search.go -package=mock joblessness/haha/search/interfaces SearchUseCase

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/auth/usecase/mock"
	"joblessness/haha/middleware"
	"joblessness/haha/models"
	searchInterfaces "joblessness/haha/search/interfaces"
	searchUseCaseMock "joblessness/haha/search/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type userSuite struct {
	suite.Suite
	router *mux.Router
	mainMiddleware *middleware.RecoveryHandler
	authMiddleware *middleware.SessionHandler
	controller *gomock.Controller
	authUseCase *mock.MockAuthUseCase
	uc *searchUseCaseMock.MockSearchUseCase
	person models.Person
	personByte *bytes.Buffer
	organization models.Organization
	organizationByte *bytes.Buffer
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/api").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.router.Use(suite.mainMiddleware.LogMiddleware)

	suite.controller = gomock.NewController(suite.T())
	suite.uc = searchUseCaseMock.NewMockSearchUseCase(suite.controller)
	suite.authUseCase = mock.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.authUseCase)

	RegisterHTTPEndpoints(suite.router, suite.uc)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestSearch() {
	suite.uc.EXPECT().
		Search(gomock.Any(), gomock.Any(), "1", "true").
		Return(models.SearchResult{}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/search?type=type&request=request&since=1&desc=true", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestSearchWrongReq() {
	suite.uc.EXPECT().
		Search(gomock.Any(), gomock.Any(), "1", "true").
		Return(models.SearchResult{}, searchInterfaces.ErrUnknownRequest).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/search?type=type&request=request&since=1&desc=true", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestSearchFailed() {
	suite.uc.EXPECT().
		Search(gomock.Any(), gomock.Any(), "1", "true").
		Return(models.SearchResult{}, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/search?type=type&request=request&since=1&desc=true", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}