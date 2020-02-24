package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()

	body := bytes.NewReader([]byte(`{"login": "marat1k", "password": "ABCDE12345"}`))

	r := httptest.NewRequest("POST", "/api/users/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if w.Result().Cookies()[0].Name == "session-id" {
		t.Error("Cookie wasnt received")
	}

	if len(h.sessions) != 1 {
		t.Error("Cookie wasnt saved")
	}
}

func TestFailedLoginNotFound(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()

	body := bytes.NewReader([]byte(`{"login": "maratk", "password": "ABE12345"}`))

	r := httptest.NewRequest("POST", "/api/users/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("status is not 404")
	}

	if len(w.Result().Cookies()) != 0 {
		t.Error("Wrong Cookie was received")
	}

	if len(h.sessions) == 1 {
		t.Error("Wrong Cookie wasnt saved")
	}
}

func TestFailedLoginWrongPassword(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()

	body := bytes.NewReader([]byte(`{"login": "marat1k", "password": "ABE12345"}`))

	r := httptest.NewRequest("POST", "/api/users/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not 400")
	}

	if len(w.Result().Cookies()) != 0 {
		t.Error("Wrong Cookie was received")
	}

	if len(h.sessions) == 1 {
		t.Error("Wrong Cookie wasnt saved")
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/users/logout", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "marat1k",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	w := httptest.NewRecorder()

	h.Logout(w, r)

	if w.Code != http.StatusOK{
		t.Error("status is not 200")
	}

	if len(h.sessions) != 0 {
		t.Error("Session wasnt closed")
	}
}

func TestLogoutWrongCookie(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/users/logout", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "mart1k",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	w := httptest.NewRecorder()

	h.Logout(w, r)

	if w.Code != http.StatusUnauthorized{
		t.Error("status is not 401")
	}

	if len(h.sessions) == 0 {
		t.Error("Wrong session was closed")
	}
}

func TestLogoutNoCookie(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/users/logout", body)
	w := httptest.NewRecorder()

	h.Logout(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("status is not 401")
	}

	if len(h.sessions) == 0 {
		t.Error("Wrong session was closed")
	}
}

func TestRegistration(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()

	//TODO Не забывать пробелы!!!
	body := bytes.NewReader([]byte(`{"login": "huvalk", 
"password": "ABE12345", 
"first-name": "first", 
"last-name": "last", 
"email": "m@m.m", 
"phone-number": "89032909812"
}`))

	r := httptest.NewRequest("POST", "/api/users", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not 200")
	}

	expectedUser := User{
			ID: 2,
			Login: "huvalk",
			Password: "ABE12345",
			FirstName: "first",
			LastName: "last",
			Email: "m@m.m",
			PhoneNumber: "89032909812",
		}

	reflect.DeepEqual(h.users["huvalk"], expectedUser)
}

func TestFailedRegistration(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()

	//TODO Не забывать пробелы!!!
	body := bytes.NewReader([]byte(`{"login": "marat1k", 
"password": "ABE12345", 
"first-name": "first", 
"last-name": "last", 
"email": "m@m.m", 
"phone-number": "89032909812"
}`))

	r := httptest.NewRequest("POST", "/api/users", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not 400")
	}

	if len(h.users) != 1 {
		t.Error("Wrong user was created")
	}
}

func TestGetUserPage(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "marat1k",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})
	w := httptest.NewRecorder()

	h.GetUserPage(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not 200")
	}

	//bytes, _ := ioutil.ReadAll(w.Body)
	//expectedJSON := `{"user":{"first-name":"Marat","last-name":"Ishimbaev","avatar":"default-avatar.jpg"},"summaries":[]}`
	//if string(bytes) != expectedJSON {
	//	t.Errorf("expected: [%s],\n got: [%s]", expectedJSON, string(bytes))
	//}
}

//TODO можно добавить еще тестов
func TestFailedGetUserPageNoUserFound(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "marat1k",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "2"})
	w := httptest.NewRecorder()

	h.GetUserPage(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("status is not 404")
	}
}

func TestSetUserInfo(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1

	body := bytes.NewReader([]byte(`{"first-name": "maratk", 
"last-name": "last", 
"password": "ABCDE2345"}`))

	r := httptest.NewRequest("PUT", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "marat1k",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.SetUserInfo(w, r)

	if w.Code != http.StatusOK{
		t.Error("status is not 200")
	}

	if (*h.users["marat1k"]).FirstName != "maratk" {
		t.Error("Changes werent saved")
	}
}

//TODO можно добавить еще тестов
func TestFailedSetUserInfoNoRights(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["maratk"] = 1
	h.sessions["marat1k"] = 2

	body := bytes.NewReader([]byte(`{"first-name": "maratk", 
"last-name": "last", 
"password": "ABCDE2345"}`))

	r := httptest.NewRequest("PUT", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "marat1k",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.SetUserInfo(w, r)

	if w.Code != http.StatusForbidden {
		t.Error("status is not 403")
	}
}

func TestSetAvatar(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1

	body := bytes.NewReader([]byte(`{"avatar": "avatar"}`))

	r := httptest.NewRequest("POST", "/api/users/1/avatar", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "marat1k",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.SetAvatar(w, r)

	if w.Code != http.StatusOK{
		t.Error("status is not 200")
	}

	if h.userAvatars[1] != "avatar" {
		t.Error("Changes werent saved")
	}
}

func TestFaildSetAvatarNoRights(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler()
	h.sessions["marat1k"] = 1
	h.sessions["maratk"] = 2
	h.userAvatars[1] = "before"

	body := bytes.NewReader([]byte(`{"avatar": "avatar"}`))

	r := httptest.NewRequest("POST", "/api/users/1/avatar", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "maratk",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.SetAvatar(w, r)

	if w.Code != http.StatusForbidden{
		t.Error("status is not 403")
	}

	if h.userAvatars[1] != "before" {
		t.Error("Wrong Changes were saved")
	}
}
