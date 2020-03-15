package httpAuth

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"joblessness/haha/auth"
	"joblessness/haha/auth/usecase/mock"
	"joblessness/haha/middleware"
	"joblessness/haha/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegistration(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	person := models.Person{
		Login:       "new username",
		Password:    "NewPassword123",
		FirstName:   "new first name",
		LastName:    "new last name",
		Email:       "new email",
		PhoneNumber: "new phone number",
	}
	personJSON, err := json.Marshal(person)
	assert.NoError(t, err)

	uc.EXPECT().
		RegisterPerson(person).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(personJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code, "Status is not 201")
}

func TestFailedRegistration(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	person := models.Person{
		Login:       "new username",
		Password:    "NewPassword123",
		FirstName:   "new first name",
		LastName:    "new last name",
		Email:       "new email",
		PhoneNumber: "new phone number",
	}
	personJSON, err := json.Marshal(person)
	assert.NoError(t, err)

	uc.EXPECT().
		RegisterPerson(&person).
		Return(auth.ErrUserAlreadyExists).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(personJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code, "Status is not 400")
}

func TestLogin(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	userLogin := models.UserLogin{
		Login:    "username",
		Password: "Password123",
	}
	userJSON, err := json.Marshal(userLogin)
	assert.NoError(t, err)

	uc.EXPECT().
		Login(userLogin.Login, userLogin.Password).
		Return(uint64(1), "sid", nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(userJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code, "Status is not 201")
	assert.Equal(t, w.Result().Cookies()[0].Name, "session_id", "Cookie wasn't received")
}

func TestFailedLoginNotFound(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	userLogin := models.UserLogin{
		Login:    "username",
		Password: "Password123",
	}
	userJSON, err := json.Marshal(userLogin)
	assert.NoError(t, err)

	uc.EXPECT().
		Login(userLogin.Login, userLogin.Password).
		Return(uint64(0), "", auth.ErrWrongLogPas).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(userJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code, "Status is not 400")
}

func TestLogout(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	uc.EXPECT().
		Logout(cookie.Value).
		Return(nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/logout", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code, "Status is not 201")
}

func TestLogoutNoCookie(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	uc.EXPECT().
		Logout(gomock.Any()).
		Times(0)

	r, _ := http.NewRequest("POST", "/api/users/logout", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 401, w.Code, "Status is not 401")
}

func TestLogoutSomethingWentWrong(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	uc.EXPECT().
		Logout(gomock.Any()).
		Return(errors.New("err")).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/logout", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 500, w.Code, "Status is not 500")
}

func TestCheck(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(1), nil).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code, "Status is not 201")
}

func TestCheckNoCookie(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)


	uc.EXPECT().
		SessionExists(gomock.Any()).
		Times(0)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 401, w.Code, "Status is not 401")
}

func TestCheckWrongSid(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(0), auth.ErrWrongSID).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 401, w.Code, "Status is not 401")
}

func TestCheckSomethingWentWrong(t *testing.T) {
	t.Parallel()
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	controller := gomock.NewController(t)
	uc := mock.NewMockUseCase(controller)
	m := middleware.NewAuthMiddleware(uc)

	RegisterHTTPEndpoints(router, m, uc)

	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}

	uc.EXPECT().
		SessionExists(cookie.Value).
		Return(uint64(0), errors.New("err")).
		Times(1)

	r, _ := http.NewRequest("POST", "/api/users/check", bytes.NewBuffer([]byte{}))
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 500, w.Code, "Status is not 500")
}