package authHttp

//go:generate mockgen -destination=../../usecase/mock/usecase.go -package=mock joblessness/haha/auth/interfaces AuthUseCase

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/status"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/auth/usecase/mock"
	"joblessness/haha/middleware"
	"joblessness/haha/models/base"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	router           *mux.Router
	mainMiddleware   *middleware.RecoveryHandler
	authMiddleware   *middleware.SessionHandler
	controller       *gomock.Controller
	uc               *mock.MockAuthUseCase
	person           baseModels.Person
	personByte       *bytes.Buffer
	organization     baseModels.Organization
	organizationByte *bytes.Buffer
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/haha").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.router.Use(suite.mainMiddleware.LogMiddleware)

	suite.controller = gomock.NewController(suite.T())
	suite.uc = mock.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.uc)

	suite.person = baseModels.Person{
		ID:        12,
		Login:     "new username",
		Password:  "NewPassword123",
		FirstName: "new first name",
		LastName:  "new last name",
		Email:     "new@email.ru",
		Phone:     "new phone number",
	}
	var err error
	personJSON, err := json.Marshal(suite.person)
	suite.personByte = bytes.NewBuffer(personJSON)
	assert.NoError(suite.T(), err)

	suite.organization = baseModels.Organization{
		ID:       12,
		Login:    "new username",
		Password: "NewPassword123",
		Name:     "new name",
		Site:     "new site",
		Email:    "new@email.ru",
		Phone:    "new phone number",
		Tag:      "awdawdawd",
	}
	organizationJSON, err := json.Marshal(suite.organization)
	suite.organizationByte = bytes.NewBuffer(organizationJSON)
	assert.NoError(suite.T(), err)

	RegisterHTTPEndpoints(suite.router, suite.authMiddleware, suite.uc)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

const message = `
--MyBoundary
Content-Disposition: form-data; name="file"; filename="file.png"
Content-Type: text/plain
`

func newTestMultipartRequest(t *testing.T) *http.Request {
	b := strings.NewReader(strings.ReplaceAll(message, "\n", "\r\n"))
	req, err := http.NewRequest("POST", "/haha/users/12/avatar", b)
	if err != nil {
		t.Fatal("NewRequest:", err)
	}
	ctype := `multipart/form-data; boundary="MyBoundary"`
	req.Header.Set("Content-type", ctype)
	return req
}

func (suite *userSuite) TestRegistrationPerson() {

	suite.uc.EXPECT().
		RegisterPerson(suite.person.Login, suite.person.Password, suite.person.FirstName).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users", suite.personByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestFailedRegistrationPerson() {
	suite.uc.EXPECT().
		RegisterPerson(suite.person.Login, suite.person.Password, suite.person.FirstName).
		Return(status.Error(authInterfaces.AlreadyExists, "user already exists")).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users", suite.personByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestRegistrationOrganization() {
	suite.uc.EXPECT().
		RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/organizations", suite.organizationByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestFailedRegistrationOrganization() {
	suite.uc.EXPECT().
		RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name).
		Return(status.Error(authInterfaces.AlreadyExists, "user already exists")).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/organizations", suite.organizationByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestLogin() {

	userLogin := baseModels.UserLogin{
		Login:    "username",
		Password: "Password123",
	}
	userJSON, err := json.Marshal(userLogin)
	assert.NoError(suite.T(), err)

	suite.uc.EXPECT().
		Login(userLogin.Login, userLogin.Password).
		Return(uint64(1), "organization", "sid", nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users/login", bytes.NewBuffer(userJSON))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
	assert.Equal(suite.T(), w.Result().Cookies()[0].Name, "session_id", "Cookie wasn't received")
}

func (suite *userSuite) TestFailedLoginNotFound() {
	userLogin := baseModels.UserLogin{
		Login:    "username",
		Password: "Password123",
	}
	userJSON, err := json.Marshal(userLogin)
	assert.NoError(suite.T(), err)

	suite.uc.EXPECT().
		Login(userLogin.Login, userLogin.Password).
		Return(uint64(0), "organization", "", status.Error(authInterfaces.WrongLoginOrPassword, "wrong login or password")).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users/login", bytes.NewBuffer(userJSON))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestLogout() {

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		Logout(cookie.Value).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users/logout", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestLogoutNoCookie() {

	suite.uc.EXPECT().
		Logout(gomock.Any()).
		Times(0)

	r, _ := http.NewRequest("POST", "/haha/users/logout", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestLogoutSomethingWentWrong() {

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		Logout(gomock.Any()).
		Return(errors.New("err")).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users/logout", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestCheck() {

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(1), nil).
		Times(1)
	suite.uc.EXPECT().
		GetRole(uint64(1)).
		Return("organization", nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestCheckNoCookie() {
	suite.uc.EXPECT().
		SessionExists(gomock.Any()).
		Times(0)

	r, _ := http.NewRequest("POST", "/haha/users/check", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestCheckWrongSid() {

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(0), authInterfaces.ErrWrongSID).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestCheckSomethingWentWrong() {

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(0), errors.New("err")).
		Times(1)

	r, _ := http.NewRequest("POST", "/haha/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}
