package tests

import (
	_h "../handlers"
	_models "../models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

func NewNotEmptyAuthHandler() *_h.AuthHandler {
	return &_h.AuthHandler {
		Sessions: make(map[string]uint, 10),
		Users:    map[string]*_models.User {
			"username": {1, "username", "Password123", "first name", "last name", "email", "phone number"},
		},
		UserAvatars: map[uint]string{},
		UserSummary: map[uint]_models.UserSummary{},
		Mu:          sync.RWMutex{},
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()

	userLogin, _ := json.Marshal(_models.UserLogin{
		Login:    "username",
		Password: "Password123",
	})

	body := bytes.NewReader(userLogin)

	r := httptest.NewRequest("POST", "/api/users/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusCreated {
		t.Error("Status is not 201")
	}

	if w.Result().Cookies()[0].Name == "session-id" {
		t.Error("Cookie wasn't received")
	}

	if len(h.Sessions) != 1 {
		t.Error("Cookie wasn't saved")
	}
}

func TestFailedLoginNotFound(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()

	userLogin, _ := json.Marshal(_models.UserLogin{
		Login:    "wrong username",
		Password: "Password123",
	})

	body := bytes.NewReader(userLogin)

	r := httptest.NewRequest("POST", "/api/users/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}

	if len(w.Result().Cookies()) != 0 {
		t.Error("Wrong cookie was received")
	}

	if len(h.Sessions) == 1 {
		t.Error("Wrong cookie wasn't saved")
	}
}

func TestFailedLoginWrongPassword(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()

	userLogin, _ := json.Marshal(_models.UserLogin{
		Login:    "username",
		Password: "WrongPassword123",
	})

	body := bytes.NewReader(userLogin)

	r := httptest.NewRequest("POST", "/api/users/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	if len(w.Result().Cookies()) != 0 {
		t.Error("Wrong cookie was received")
	}

	if len(h.Sessions) == 1 {
		t.Error("Wrong cookie wasn't saved")
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()
	h.Sessions["username"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/users/logout", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	w := httptest.NewRecorder()

	h.Logout(w, r)

	if w.Code != http.StatusCreated{
		t.Error("Status is not 201")
	}

	if len(h.Sessions) != 0 {
		t.Error("Session wasn't closed")
	}
}

func TestLogoutWrongCookie(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()
	h.Sessions["username"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/users/logout", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "wrong username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	w := httptest.NewRecorder()

	h.Logout(w, r)

	if w.Code != http.StatusUnauthorized{
		t.Error("Status is not 401")
	}

	if len(h.Sessions) == 0 {
		t.Error("Wrong session was closed")
	}
}

func TestLogoutNoCookie(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()
	h.Sessions["username"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/users/logout", body)
	w := httptest.NewRecorder()

	h.Logout(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("Status is not 401")
	}

	if len(h.Sessions) == 0 {
		t.Error("Wrong session was closed")
	}
}

func TestRegistration(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()

	newUser, _ := json.Marshal(_models.User{
		Login:       "new username",
		Password:    "NewPassword123",
		FirstName:   "new first name",
		LastName:    "new last name",
		Email:       "new email",
		PhoneNumber: "new phone number",
	})

	body := bytes.NewReader(newUser)

	r := httptest.NewRequest("POST", "/api/users", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusCreated {
		t.Error("Status is not 201")
	}

	expectedUser := _models.User{
		ID: 2,
		Login: "new username",
		Password: "NewPassword123",
		FirstName: "new first name",
		LastName: "new last name",
		Email: "new email",
		PhoneNumber: "new phone number",
	}

	reflect.DeepEqual(h.Users["new username"], expectedUser)
}

func TestFailedRegistration(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyAuthHandler()

	newUser, _ := json.Marshal(_models.User{
		Login:       "username",
		Password:    "Password123",
		FirstName:   "first name",
		LastName:    "last name",
		Email:       "email",
		PhoneNumber: "phone number",
	})

	body := bytes.NewReader(newUser)

	r := httptest.NewRequest("POST", "/api/users", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	if len(h.Users) != 1 {
		t.Error("Wrong user was created")
	}
}
