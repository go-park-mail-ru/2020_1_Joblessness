package userHttp

//go:generate mockgen -destination=../../usecase/mock/user.go -package=mock joblessness/haha/user/interfaces UserUseCase

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
	"joblessness/haha/models/base"
	"joblessness/haha/user/usecase/mock"
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
	uc               *mock.MockUserUseCase
	ucAuth           *mockAuth.MockAuthUseCase
	person           baseModels.Person
	personByte       *bytes.Buffer
	organization     baseModels.Organization
	organizationByte *bytes.Buffer
}

func (suite *userSuite) SetupTest() {
	suite.router = mux.NewRouter().PathPrefix("/api").Subrouter()
	suite.mainMiddleware = middleware.NewMiddleware()
	suite.router.Use(suite.mainMiddleware.LogMiddleware)

	suite.controller = gomock.NewController(suite.T())
	suite.uc = mock.NewMockUserUseCase(suite.controller)
	suite.ucAuth = mockAuth.NewMockAuthUseCase(suite.controller)
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.ucAuth)

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
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
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
	suite.ucAuth.EXPECT().
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
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(13), nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
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
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("POST", "/api/users/12/avatar", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 416, w.Code, "Status is not 416")
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
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangePerson(&suite.person).
		Return(nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("PUT", "/api/users/12", suite.personByte)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 500")
}

func (suite *userSuite) TestChangePersonNoCookie() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangePerson(&suite.person).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("PUT", "/api/users/12", suite.personByte)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestChangePersonWrongId() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangePerson(&suite.person).
		Return(nil).
		Times(0)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
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
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangeOrganization(&suite.organization).
		Return(nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("PUT", "/api/organizations/12", suite.organizationByte)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 204, w.Code, "Status is not 204")
}

func (suite *userSuite) TestChangeOrganizationNoCookie() {
	suite.ucAuth.EXPECT().
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
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		ChangeOrganization(suite.organization).
		Return(nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
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
		Return(baseModels.Organizations{}, nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations?page=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestListOrgsFailed() {
	suite.uc.EXPECT().
		GetListOfOrgs(1).
		Return(nil, errors.New("")).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/organizations?page=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestLikeUser() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		LikeUser(uint64(12), uint64(1)).
		Return(false, nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("POST", "/api/users/1/like", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestLikeUserNoSession() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/1/like", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestLikeUserFailed() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		LikeUser(uint64(12), uint64(1)).
		Return(false, errors.New("")).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("POST", "/api/users/1/like", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestLikeExists() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		LikeExists(uint64(12), uint64(1)).
		Return(false, nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("GET", "/api/users/1/like", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestLikeExistsFailed() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		LikeExists(uint64(12), uint64(1)).
		Return(false, errors.New("")).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("GET", "/api/users/1/like", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}

func (suite *userSuite) TestGetFavorite() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		GetUserFavorite(uint64(12)).
		Return(baseModels.Favorites{}, nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("GET", "/api/users/12/favorite", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 200, w.Code, "Status is not 200")
}

func (suite *userSuite) TestGetFavoriteNoSession() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	r, _ := http.NewRequest("GET", "/api/users/12/favorite", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 401, w.Code, "Status is not 401")
}

func (suite *userSuite) TestGetFavoriteWrongId() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("GET", "/api/users/13/favorite", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 403, w.Code, "Status is not 403")
}

func (suite *userSuite) TestGetFavoriteFailed() {
	suite.ucAuth.EXPECT().
		SessionExists("username").
		Return(uint64(12), nil).
		Times(1)
	suite.uc.EXPECT().
		GetUserFavorite(uint64(12)).
		Return(nil, errors.New("")).
		Times(1)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "username",
		Expires: time.Now().Add(time.Hour),
	}

	r, _ := http.NewRequest("GET", "/api/users/12/favorite", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), 500, w.Code, "Status is not 500")
}
