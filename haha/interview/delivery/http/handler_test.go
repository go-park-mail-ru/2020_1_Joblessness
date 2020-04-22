package interviewHttp

//go:generate mockgen -destination=../../usecase/mock/usecase.go -package=mock joblessness/haha/interview/interfaces InterviewUseCase

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mockAuth "joblessness/haha/auth/usecase/mock"
	interviewUseCaseMock "joblessness/haha/interview/usecase/mock"
	"joblessness/haha/middleware"
	"joblessness/haha/models/base"
	summaryInterfaces "joblessness/haha/summary/interfaces"
	"joblessness/haha/utils/chat/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	room           *mock.MockRoom
	router         *mux.Router
	mainMiddleware *middleware.RecoveryHandler
	authMiddleware *middleware.SessionHandler
	controller     *gomock.Controller
	authUseCase    *mockAuth.MockAuthUseCase
	uc             *interviewUseCaseMock.MockInterviewUseCase
	cookie         *http.Cookie
	response       baseModels.VacancyResponse
	sendSum        baseModels.SendSummary
	params         baseModels.ChatParameters
	summaryCredits baseModels.SummaryCredentials
	responseByte   *bytes.Buffer
	sendSumByte    *bytes.Buffer
	paramsByte     *bytes.Buffer
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/api").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.router.Use(suite.mainMiddleware.LogMiddleware)

	suite.controller = gomock.NewController(suite.T())
	suite.uc = interviewUseCaseMock.NewMockInterviewUseCase(suite.controller)
	suite.authUseCase = mockAuth.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.authUseCase)

	suite.cookie = &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.sendSum = baseModels.SendSummary{
		VacancyID:      uint64(7),
		SummaryID:      uint64(1),
		UserID:         uint64(2),
		OrganizationID: uint64(13),
		Accepted:       true,
		Denied:         false,
	}
	sendSumJSON, err := json.Marshal(suite.sendSum)
	suite.sendSumByte = bytes.NewBuffer(sendSumJSON)
	assert.NoError(suite.T(), err)

	suite.params = baseModels.ChatParameters{
		From: 1,
		To:   1,
		Page: 0,
	}
	paramsJSON, err := json.Marshal(suite.sendSum)
	suite.paramsByte = bytes.NewBuffer(paramsJSON)
	assert.NoError(suite.T(), err)

	suite.room = mock.NewMockRoom(suite.controller)
	suite.room.EXPECT().Run().Times(1)

	suite.summaryCredits = baseModels.SummaryCredentials{}

	RegisterHTTPEndpoints(suite.router, suite.authMiddleware, suite.uc)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestResponseSummary() {
	suite.room.EXPECT().
		SendGeneratedMessage(gomock.Any()).
		Return().
		Times(1)
	suite.uc.EXPECT().
		ResponseSummary(&suite.sendSum).
		Return(nil).
		Times(1)
	suite.uc.EXPECT().
		GetResponseCredentials(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		Return(&suite.summaryCredits, nil).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/1/response", suite.sendSumByte)
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

	r, _ := http.NewRequest("PUT", "/api/summaries/1/response", bytes.NewBuffer([]byte{}))
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
		Return(summaryInterfaces.ErrOrganizationIsNotOwner).
		Times(1)
	suite.authUseCase.EXPECT().
		OrganizationSession("username").
		Return(uint64(13), nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/summaries/1/response", suite.sendSumByte)
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

	r, _ := http.NewRequest("PUT", "/api/summaries/1/response", suite.sendSumByte)
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

	r, _ := http.NewRequest("PUT", "/api/summaries/1/response", suite.sendSumByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestHistory() {
	suite.uc.EXPECT().
		GetHistory(&suite.params).
		Return(baseModels.Messages{}, nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/chat/conversation/1", suite.paramsByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestHistoryFailed() {
	suite.uc.EXPECT().
		GetHistory(&suite.params).
		Return(baseModels.Messages{}, errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/chat/conversation/1", suite.paramsByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetConversation() {
	suite.uc.EXPECT().
		GetConversations(suite.params.From).
		Return(baseModels.Conversations{}, nil).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/chat/conversation", suite.paramsByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetConversationFailed() {
	suite.uc.EXPECT().
		GetConversations(suite.params.From).
		Return(baseModels.Conversations{}, errors.New("")).
		Times(1)
	suite.authUseCase.EXPECT().
		SessionExists("username").
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/chat/conversation", suite.paramsByte)
	r.AddCookie(suite.cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}
