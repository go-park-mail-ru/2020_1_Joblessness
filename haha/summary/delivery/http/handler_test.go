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
	"joblessness/haha/middleware/xss"
	"joblessness/haha/models"
	summaryInterfaces "joblessness/haha/summary/interfaces"
	summaryUseCaseMock "joblessness/haha/summary/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	router *mux.Router
	mainMiddleware *middleware.RecoveryHandler
	authMiddleware *middleware.SessionHandler
	xssMiddleware *xss.XssHandler
	controller *gomock.Controller
	authUseCase *mockAuth.MockAuthUseCase
	uc *summaryUseCaseMock.MockSummaryUseCase
	summary models.Summary
	summaryByte *bytes.Buffer
	cookie *http.Cookie
	response models.VacancyResponse
	sendSum models.SendSummary
	responseByte *bytes.Buffer
	sendSumByte *bytes.Buffer
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/api").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.xssMiddleware = xss.NewXssHandler()
	suite.router.Use(suite.xssMiddleware.SanitizeMiddleware)
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
		},
		Keywords:    "key",
		Educations:  []models.Education{
			models.Education{
				Institution: "was",
				Speciality:  "is",
				Type:        "none",
			},
		},
		Experiences: []models.Experience{
			models.Experience{
				CompanyName:      "comp",
				Role:             "role",
				Responsibilities: "response",
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

	suite.response = models.VacancyResponse{
		UserID:    suite.summary.Author.ID,
		Tag:       suite.summary.Author.Tag,
		VacancyID: uint64(7),
		SummaryID: suite.summary.ID,
		Keywords:  suite.summary.Keywords,
	}
	responseJSON, err := json.Marshal(suite.response)
	suite.responseByte = bytes.NewBuffer(responseJSON)
	assert.NoError(suite.T(), err)

	suite.sendSum = models.SendSummary{
		VacancyID:      uint64(7),
		SummaryID:      suite.summary.ID,
		UserID:    		suite.summary.Author.ID,
		OrganizationID: uint64(13),
		Accepted:       true,
		Denied:         false,
	}
	sendSumJSON, err := json.Marshal(suite.sendSum)
	suite.sendSumByte = bytes.NewBuffer(sendSumJSON)
	assert.NoError(suite.T(), err)

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
		PersonSession("username").
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
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/summaries", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestCreateSummaryFailed() {
	suite.uc.EXPECT().
		CreateSummary(&suite.summary).
		Return(uint64(0), errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
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

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestGetSummaries() {
	suite.uc.EXPECT().
		GetAllSummaries("").
		Return([]models.Summary{suite.summary}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetSummariesFailed() {
	suite.uc.EXPECT().
		GetAllSummaries("1").
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/summaries?page=1", bytes.NewBuffer([]byte{}))
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

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestGetUserSummaries() {
	suite.uc.EXPECT().
		GetUserSummaries(suite.summary.Author.ID).
		Return([]models.Summary{suite.summary}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/12/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetUserSummariesWrongUrl() {
	suite.uc.EXPECT().
		GetUserSummaries(suite.summary.Author.ID).
		Return([]models.Summary{suite.summary}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/a/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestGetUserSummariesFailed() {
	suite.uc.EXPECT().
		GetUserSummaries(suite.summary.Author.ID).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/12/summaries", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangeSummary() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestChangeSummaryWrongUrl() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/a", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestChangeSummaryWrongJson() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestChangeSummaryFailed() {
	suite.uc.EXPECT().
		ChangeSummary(&suite.summary).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3", suite.summaryByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestDeleteSummary() {
	suite.uc.EXPECT().
		DeleteSummary(uint64(3)).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestDeleteSummaryWrongUrl() {
	suite.uc.EXPECT().
		DeleteSummary(uint64(3)).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/summaries/a", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestDeleteSummaryFailed() {
	suite.uc.EXPECT().
		DeleteSummary(uint64(3)).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("DELETE", "/api/summaries/3", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestSendSummary() {
	suite.uc.EXPECT().
		SendSummary(&suite.sendSum).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies/7/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestSendSummaryNotOwner() {
	suite.uc.EXPECT().
		SendSummary(&suite.sendSum).
		Return(summaryInterfaces.ErrPersonIsNotOwner).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies/7/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 403, w.Code, "Status is not 403")
}

func (suite *userSuite) TestSendSummaryNoSummary() {
	suite.uc.EXPECT().
		SendSummary(&suite.sendSum).
		Return(summaryInterfaces.ErrNoSummaryToRefresh).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies/7/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestSendSummaryDefaultErr() {
	suite.uc.EXPECT().
		SendSummary(&suite.sendSum).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies/7/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestSendSummaryWrongJson() {
	suite.uc.EXPECT().
		SendSummary(&suite.sendSum).
		Return(summaryInterfaces.ErrNoSummaryToRefresh).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies/7/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestSendSummaryWrongUrl() {
	suite.uc.EXPECT().
		SendSummary(&suite.sendSum).
		Return(summaryInterfaces.ErrNoSummaryToRefresh).
		Times(1)
	suite.authUseCase.EXPECT().
		PersonSession("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/vacancies/7a/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestResponseSummary() {
	suite.uc.EXPECT().
		ResponseSummary(&suite.sendSum).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestResponseSummaryWrongJSON() {
	suite.uc.EXPECT().
		ResponseSummary(&suite.sendSum).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestResponseSummaryWrongURL() {
	suite.uc.EXPECT().
		ResponseSummary(&suite.sendSum).
		Return(nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3a/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestResponseSummaryNotOwner() {
	suite.uc.EXPECT().
		ResponseSummary(&suite.sendSum).
		Return(summaryInterfaces.ErrOrgIsNotOwner).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 403, w.Code, "Status is not 403")
}

func (suite *userSuite) TestResponseSummaryNoSummary() {
	suite.uc.EXPECT().
		ResponseSummary(&suite.sendSum).
		Return(summaryInterfaces.ErrNoSummaryToRefresh).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestResponseSummaryDefaultErr() {
	suite.uc.EXPECT().
		ResponseSummary(&suite.sendSum).
		Return(errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/3/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetOrgSendSummaries() {
	suite.uc.EXPECT().
		GetOrgSendSummaries(suite.sendSum.UserID).
		Return([]*models.VacancyResponse{&suite.response}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/12/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetOrgSendSummariesWrongURL() {
	suite.uc.EXPECT().
		GetOrgSendSummaries(suite.sendSum.UserID).
		Times(0)

	r, _ := http.NewRequest("GET", "/api/organizations/12a/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestGetOrgSendSummariesFailed() {
	suite.uc.EXPECT().
		GetOrgSendSummaries(suite.sendSum.UserID).
		Return([]*models.VacancyResponse{}, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/12/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetUserSendSummaries() {
	suite.uc.EXPECT().
		GetUserSendSummaries(suite.sendSum.UserID).
		Return([]*models.VacancyResponse{&suite.response}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/12/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetUserSendSummariesWrongURL() {
	suite.uc.EXPECT().
		GetUserSendSummaries(suite.sendSum.UserID).
		Times(0)

	r, _ := http.NewRequest("GET", "/api/users/12a/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 404, w.Code, "Status is not 404")
}

func (suite *userSuite) TestGetUserSendSummariesFailed() {
	suite.uc.EXPECT().
		GetUserSendSummaries(suite.sendSum.UserID).
		Return([]*models.VacancyResponse{}, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/12/response", bytes.NewBuffer([]byte{}))
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}