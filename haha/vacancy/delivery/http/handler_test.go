package httpVacancy

//go:generate  mockgen -destination=../../usecase/mock/vacancy.go -package=mock joblessness/haha/vacancy/interfaces VacancyUseCase

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mockAuth "joblessness/haha/auth/usecase/mock"
	"joblessness/haha/middleware"
	"joblessness/haha/models"
	vacancyUseCaseMock "joblessness/haha/vacancy/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	router         *mux.Router
	mainMiddleware *middleware.RecoveryHandler
	authMiddleware *middleware.SessionHandler
	controller     *gomock.Controller
	authUseCase    *mockAuth.MockAuthUseCase
	uc             *vacancyUseCaseMock.MockVacancyUseCase
	vacancy        models.Vacancy
	vacancyByte    *bytes.Buffer
	cookie         *http.Cookie
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/api").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.router.Use(suite.mainMiddleware.LogMiddleware)

	suite.controller = gomock.NewController(suite.T())
	suite.uc = vacancyUseCaseMock.NewMockVacancyUseCase(suite.controller)
	suite.authUseCase = mockAuth.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.authUseCase)

	suite.vacancy = models.Vacancy{
		ID: 3,
		Organization: models.VacancyOrganization{
			ID:     12,
			Tag:    "",
			Email:  "",
			Phone:  "",
			Avatar: "",
			Name:   "",
			Site:   "",
		},
		Name:             "vacancy",
		Description:      "description",
		SalaryFrom:       50,
		SalaryTo:         100,
		WithTax:          false,
		Responsibilities: "all",
		Conditions:       "some",
		Keywords:         "word",
	}
	var err error
	vacancyJSON, err := json.Marshal(suite.vacancy)
	suite.vacancyByte = bytes.NewBuffer(vacancyJSON)
	assert.NoError(suite.T(), err)

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

func (suite *userSuite) TestCreateVacancy() {
	suite.uc.EXPECT().
		CreateVacancy(&suite.vacancy).
		Return(uint64(3), nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies", suite.vacancyByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestCreateVacancyEmptyName() {
	vacancy := suite.vacancy
	vacancy.Name = ""
	var err error
	vacancyJSON, err := json.Marshal(vacancy)
	vacancyByte := bytes.NewBuffer(vacancyJSON)
	assert.NoError(suite.T(), err)

	suite.uc.EXPECT().
		CreateVacancy(&suite.vacancy).
		Return(uint64(3), nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies", vacancyByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestCreateVacancyWrongJson() {
	suite.uc.EXPECT().
		CreateVacancy(&suite.vacancy).
		Return(uint64(3), nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestCreateVacancyFailed() {
	suite.uc.EXPECT().
		CreateVacancy(&suite.vacancy).
		Return(uint64(0), errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies", suite.vacancyByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetVacancy() {
	suite.uc.EXPECT().
		GetVacancy(uint64(3)).
		Return(&suite.vacancy, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/vacancies/3", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetVacancyFailed() {
	suite.uc.EXPECT().
		GetVacancy(uint64(3)).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/vacancies/3", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetVacancyWrongUrl() {
	suite.uc.EXPECT().
		GetVacancy(uint64(3)).
		Return(&suite.vacancy, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/vacancies/a", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestGetVacancies() {
	suite.uc.EXPECT().
		GetVacancies("1").
		Return(models.Vacancies{&suite.vacancy}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/vacancies?page=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetVacanciesEmpty() {
	suite.uc.EXPECT().
		GetVacancies("1").
		Return(models.Vacancies{}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/vacancies?page=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetVacanciesFailed() {
	suite.uc.EXPECT().
		GetVacancies("1").
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/vacancies?page=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangeVacancy() {
	suite.uc.EXPECT().
		ChangeVacancy(&suite.vacancy).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/vacancies/3", suite.vacancyByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestChangeVacancyWrongUrl() {
	suite.uc.EXPECT().
		ChangeVacancy(&suite.vacancy).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/vacancies/a", suite.vacancyByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestChangeVacancyWrongJson() {
	suite.uc.EXPECT().
		ChangeVacancy(&suite.vacancy).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/vacancies/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestChangeVacancyFailed() {
	suite.uc.EXPECT().
		ChangeVacancy(&suite.vacancy).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/vacancies/3", suite.vacancyByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestDeleteVacancy() {
	suite.uc.EXPECT().
		DeleteVacancy(uint64(3), uint64(12)).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/vacancies/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestDeleteVacancyWrongUrl() {
	suite.uc.EXPECT().
		DeleteVacancy(uint64(3), uint64(12)).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/vacancies/a", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestDeleteVacancyFailed() {
	suite.uc.EXPECT().
		DeleteVacancy(uint64(3), uint64(12)).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/vacancies/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetOrgVacancies() {
	suite.uc.EXPECT().
		GetOrgVacancies(uint64(1)).
		Return(models.Vacancies{&suite.vacancy}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/1/vacancies", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetOrgVacanciesEmpty() {
	suite.uc.EXPECT().
		GetOrgVacancies(uint64(1)).
		Return(models.Vacancies{}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/1/vacancies", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetOrgVacanciesFailed() {
	suite.uc.EXPECT().
		GetOrgVacancies(uint64(1)).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/1/vacancies", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}
