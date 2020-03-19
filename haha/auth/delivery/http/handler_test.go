package httpAuth

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/auth"
	"joblessness/haha/auth/usecase/mock"
	"joblessness/haha/middleware"
	"joblessness/haha/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	router *mux.Router
	mainMiddleware *middleware.Middleware
	authMiddleware *middleware.AuthMiddleware
	controller *gomock.Controller
	uc *mock.MockAuthUseCase
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
	suite.uc = mock.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.uc)

	suite.person = models.Person{
		ID: 12,
		Login:       "new username",
		Password:    "NewPassword123",
		FirstName:   "new first name",
		LastName:    "new last name",
		Email:       "new email",
		Phone: "new phone number",
	}
	var err error
	personJSON, err := json.Marshal(suite.person)
	suite.personByte = bytes.NewBuffer(personJSON)
	assert.NoError(suite.T(), err)

	suite.organization = models.Organization{
		ID: 12,
		Login:       "new username",
		Password:    "NewPassword123",
		Name:   "new name",
		Site:    "new site",
		Email:       "new email",
		Phone: "new phone number",
	}
	organizationJSON, err := json.Marshal(suite.organization)
	suite.organizationByte = bytes.NewBuffer(organizationJSON)
	assert.NoError(suite.T(), err)

	RegisterHTTPEndpoints(suite.router,suite.authMiddleware, suite.uc)
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
	req, err := http.NewRequest("POST", "/api/users/12/avatar", b)
	if err != nil {
		t.Fatal("NewRequest:", err)
	}
	ctype := `multipart/form-data; boundary="MyBoundary"`
	req.Header.Set("Content-type", ctype)
	return req
}

func (suite *userSuite) TestSetAvatar() {
	suite.uc.EXPECT().
		SetAvatar(gomock.Any(), uint64(12)).
		Return(nil).
		Times(1)
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	r := newTestMultipartRequest(suite.T())
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestSetAvatarNoCookie() {
	suite.uc.EXPECT().
		SetAvatar(gomock.Any(), uint64(12)).
		Return(nil).
		Times(1)
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r := newTestMultipartRequest(suite.T())
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestSetAvatarWrongId() {
	suite.uc.EXPECT().
		SetAvatar(gomock.Any(), uint64(12)).
		Return(nil).
		Times(0)
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(13), nil).
		Times(1)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	r := newTestMultipartRequest(suite.T())
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 403, w.Code, "Status is not 403")
}

func (suite *userSuite) TestSetAvatarNotMultipart() {
	suite.uc.EXPECT().
		SetAvatar(gomock.Any(), uint64(12)).
		Return(nil).
		Times(1)
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("POST", "/api/users/12/avatar", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 416, w.Code, "Status is not 416")
}

func (suite *userSuite) TestRegistrationPerson() {

	suite.uc.EXPECT().
		RegisterPerson(&suite.person).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users", suite.personByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestFailedRegistrationPerson() {
	suite.uc.EXPECT().
		RegisterPerson(&suite.person).
		Return(auth.ErrUserAlreadyExists).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users", suite.personByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestRegistrationOrganization() {
	suite.uc.EXPECT().
		RegisterOrganization(&suite.organization).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/organizations", suite.organizationByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestFailedRegistrationOrganization() {
	suite.uc.EXPECT().
		RegisterOrganization(&suite.organization).
		Return(auth.ErrUserAlreadyExists).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/organizations", suite.organizationByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestLogin() {

	userLogin := models.UserLogin{
		Login:    "username",
		Password: "Password123",
	}
	userJSON, err := json.Marshal(userLogin)
	assert.NoError(suite.T(), err)

	suite.uc.EXPECT().
		Login(userLogin.Login, userLogin.Password).
		Return(uint64(1), "sid", nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(userJSON))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
	assert.Equal(suite.T(), w.Result().Cookies()[0].Name, "session_id", "Cookie wasn't received")
}

func (suite *userSuite) TestFailedLoginNotFound() {

	userLogin := models.UserLogin{
		Login:    "username",
		Password: "Password123",
	}
	userJSON, err := json.Marshal(userLogin)
	assert.NoError(suite.T(), err)

	suite.uc.EXPECT().
		Login(userLogin.Login, userLogin.Password).
		Return(uint64(0), "", auth.ErrWrongLogPas).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(userJSON))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 400, w.Code, "Status is not 400")
}

func (suite *userSuite) TestLogout() {

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		Logout(cookie.Value).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/logout", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestLogoutNoCookie() {

	suite.uc.EXPECT().
		Logout(gomock.Any()).
		Times(0)

	r, _ := http.NewRequest("POST", "/api/users/logout", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestLogoutSomethingWentWrong() {

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		Logout(gomock.Any()).
		Return(errors.New("err")).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/logout", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestCheck() {

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 201, w.Code, "Status is not 201")
}

func (suite *userSuite) TestCheckNoCookie() {


	suite.uc.EXPECT().
		SessionExists(gomock.Any()).
		Times(0)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestCheckWrongSid() {

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(0), auth.ErrWrongSID).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestCheckSomethingWentWrong() {

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	suite.uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(0), errors.New("err")).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetPerson() {
	suite.uc.EXPECT().
		GetPerson(uint64(12)).
		Return(&suite.person, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/12", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetPersonWentWrong() {
	suite.uc.EXPECT().
		GetPerson(uint64(12)).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/12", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangePerson() {
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangePerson(suite.person).
		Return(nil).
		Times(1)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("PUT", "/api/users/12", suite.personByte)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangePersonNoCookie() {
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangePerson(suite.person).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/users/12", suite.personByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestChangePersonWrongId() {
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangePerson(suite.person).
		Return(nil).
		Times(0)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("PUT", "/api/users/13", suite.personByte)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 403, w.Code, "Status is not 403")
}

func (suite *userSuite) TestGetOrganization() {
	suite.uc.EXPECT().
		GetOrganization(uint64(12)).
		Return(&suite.organization, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/12", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetOrganizationWentWrong() {
	suite.uc.EXPECT().
		GetOrganization(uint64(12)).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/12", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangeOrganization() {
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangeOrganization(suite.organization).
		Return(nil).
		Times(1)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("PUT", "/api/organizations/12", suite.organizationByte)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestChangeOrganizationNoCookie() {
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangeOrganization(suite.organization).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/organizations/12", suite.organizationByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestChangeOrganizationWrongId() {
	suite.uc.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangeOrganization(suite.organization).
		Return(nil).
		Times(1)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("PUT", "/api/organizations/13", suite.organizationByte)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 403, w.Code, "Status is not 403")
}

func (suite *userSuite) TestListOrgs() {
	suite.uc.EXPECT().
		GetListOfOrgs(1).
		Return([]models.Organization{}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/?page=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestListOrgsFailed() {
	suite.uc.EXPECT().
		GetListOfOrgs(1).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations/?page=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}