package tests

import (
	_h "../handlers"
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserPage(t *testing.T) {
	t.Parallel()

	h := _h.NewAuthHandler()
	h.Sessions["marat1k"] = 1

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

	h := _h.NewAuthHandler()
	h.Sessions["marat1k"] = 1

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

func TestChangeUserInfo(t *testing.T) {
	t.Parallel()

	h := _h.NewAuthHandler()
	h.Sessions["marat1k"] = 1

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

	h.ChangeUserInfo(w, r)

	if w.Code != http.StatusNoContent{
		t.Error("status is not 204")
	}

	if (*h.Users["marat1k"]).FirstName != "maratk" {
		t.Error("Changes werent saved")
	}
}

//TODO можно добавить еще тестов
func TestFailedChangeUserInfoNoRights(t *testing.T) {
	t.Parallel()

	h := _h.NewAuthHandler()
	h.Sessions["maratk"] = 1
	h.Sessions["marat1k"] = 2

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

	h.ChangeUserInfo(w, r)

	if w.Code != http.StatusForbidden {
		t.Error("status is not 403")
	}
}

func TestSetAvatar(t *testing.T) {
	t.Parallel()

	h := _h.NewAuthHandler()
	h.Sessions["marat1k"] = 1

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

	if w.Code != http.StatusCreated{
		t.Error("status is not 201")
	}

	if h.UserAvatars[1] != "avatar" {
		t.Error("Changes werent saved")
	}
}

func TestFaildSetAvatarNoRights(t *testing.T) {
	t.Parallel()

	h := _h.NewAuthHandler()
	h.Sessions["marat1k"] = 1
	h.Sessions["maratk"] = 2
	h.UserAvatars[1] = "before"

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

	if h.UserAvatars[1] != "before" {
		t.Error("Wrong Changes were saved")
	}
}