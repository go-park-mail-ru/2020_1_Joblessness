package httpSummary

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
	summaryUseCaseMock "joblessness/haha/summary/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	router *mux.Router
	mainMiddleware *middleware.Middleware
	authMiddleware *middleware.AuthMiddleware
	controller *gomock.Controller
	authUseCase *mockAuth.MockAuthUseCase
	uc *summaryUseCaseMock.MockSummaryUseCase
	summary models.Summary
	summaryByte *bytes.Buffer
	cookie *http.Cookie
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/api").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.router.Use(suite.mainMiddleware.LogMiddleware)

	suite.controller = gomock.NewController(suite.T())
	suite.uc = summaryUseCaseMock.NewMockSummaryUseCase(suite.controller)
	suite.authUseCase = mockAuth.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.authUseCase)

	suite.summary = models.Summary{
		ID:          3,
		Author:      models.Author{
			ID:        12,
			Tag:       "tag",
			Email:     "email",
			Phone:     "phone",
			Avatar:    "avatar",
			FirstName: "first",
			LastName:  "name",
			Gender:    "gender",
			Birthday:  time.Now(),
		},
		Keywords:    "key",
		Educations:  []models.Education{
			models.Education{
				Institution: "was",
				Speciality:  "is",
				Graduated:   time.Now(),
				Type:        "none",
			},
		},
		Experiences: []models.Experience{
			models.Experience{
				CompanyName:      "comp",
				Role:             "role",
				Responsibilities: "response",
				Start:            time.Now(),
				Stop:             time.Now().AddDate(1, 1, 1),
			},
		},
	}
	var err error
	summaryJSON, err := json.Marshal(suite.summary)
	suite.summaryByte = bytes.NewBuffer(summaryJSON)
	assert.NoError(suite.T(), err)

	suite.cookie = &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	RegisterHTTPEndpoints(suite.router,suite.authMiddleware, suite.uc)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestCreateSummary() {
	suite.uc.EXPECT().
		CreateSummary(&suite.summary).
		Return(uint64(3), nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/summaries", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestCreateSummaryWrongJson() {
	suite.uc.EXPECT().
		CreateSummary(&suite.summary).
		Return(uint64(3), nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/summaries", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestCreateSummaryFailed() {
	suite.uc.EXPECT().
		CreateSummary(&suite.summary).
		Return(uint64(0), errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/summaries", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetSummary() {
	suite.uc.EXPECT().
		GetSummary(uint64(3)).
		Return(&suite.summary, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetSummaryFailed() {
	suite.uc.EXPECT().
		GetSummary(uint64(3)).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetSummaryWrongUrl() {
	suite.uc.EXPECT().
		GetSummary(uint64(3)).
		Return(&suite.summary, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries/a", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestGetSummaries() {
	suite.uc.EXPECT().
		GetAllSummaries().
		Return([]models.Summary{suite.summary}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetSummariesFailed() {
	suite.uc.EXPECT().
		GetAllSummaries().
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestPrintSummaries() {
	suite.uc.EXPECT().
		GetSummary(suite.summary.ID).
		Return(&suite.summary, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries/3/print", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestPrintSummariesFailed() {
	suite.uc.EXPECT().
		GetSummary(suite.summary.ID).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries/3/print", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestPrintSummariesWrongUrl() {
	suite.uc.EXPECT().
		GetSummary(suite.summary.ID).
		Return(&suite.summary, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries/a/print", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestGetUserSummaries() {
	suite.uc.EXPECT().
		GetUserSummaries(suite.summary.Author.ID).
		Return([]models.Summary{suite.summary}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/user/12/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetUserSummariesWrongUrl() {
	suite.uc.EXPECT().
		GetUserSummaries(suite.summary.Author.ID).
		Return([]models.Summary{suite.summary}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/user/a/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestGetUserSummariesFailed() {
	suite.uc.EXPECT().
		GetUserSummaries(suite.summary.Author.ID).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/user/12/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangeVacancy() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestChangeVacancyWrongUrl() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/a", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestChangeVacancyWrongJson() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangeVacancyFailed() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestDeleteVacancy() {
	suite.uc.EXPECT().
		DeleteSummary(uint64(3)).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestDeleteVacancyWrongUrl() {
	suite.uc.EXPECT().
		DeleteSummary(uint64(3)).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/summaries/a", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestDeleteVacancyFailed() {
	suite.uc.EXPECT().
		DeleteSummary(uint64(3)).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}