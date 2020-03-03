package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"joblessness/haha/handlers"
	"joblessness/haha/models"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type Avatar struct {
	Avatar string
}

func NewNotEmptyUsersHandler() *handlers.AuthHandler {
	return &handlers.AuthHandler {
		Sessions: make(map[string]uint, 10),
		Users:    map[string]*models.User {
			"username": {1, "username", "Password123", "first name", "last name", "email", "phone number"},
		},
		UserAvatars: map[uint]string{},
		UserSummary: map[uint]models.UserSummary{},
		Mu:          sync.RWMutex{},
	}
}

func TestGetUserPage(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyUsersHandler()
	h.Sessions["username"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})
	w := httptest.NewRecorder()

	h.GetUserPage(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not 200")
	}
}

//TODO можно добавить еще тестов
func TestFailedGetUserPageNoUserFound(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyUsersHandler()
	h.Sessions["username"] = 1

	body := bytes.NewReader([]byte{})

	r := httptest.NewRequest("POST", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "2"})
	w := httptest.NewRecorder()

	h.GetUserPage(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}
}

func TestChangeUserInfo(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyUsersHandler()
	h.Sessions["username"] = 1

	user, _ := json.Marshal(models.User{
		Password:    "NewPassword123",
		FirstName:   "new first name",
		LastName:    "new last name",
	})

	body := bytes.NewReader(user)

	r := httptest.NewRequest("PUT", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.ChangeUserInfo(w, r)

	if w.Code != http.StatusNoContent{
		t.Error("Status is not 204")
	}

	if (*h.Users["username"]).FirstName != "new first name" {
		t.Error("Changes weren't saved")
	}
}

func TestFailedChangeUserInfoNoRights(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyUsersHandler()
	h.Sessions["other username"] = 1
	h.Sessions["username"] = 2

	user, _ := json.Marshal(models.User{
		Password:    "NewPassword123",
		FirstName:   "new first name",
		LastName:    "new last name",
	})

	body := bytes.NewReader(user)

	r := httptest.NewRequest("PUT", "/api/user", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.ChangeUserInfo(w, r)

	if w.Code != http.StatusForbidden {
		t.Error("Status is not 403")
	}
}

func TestSetAvatar(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyUsersHandler()
	h.Sessions["username"] = 1

	avatar, _ := json.Marshal(Avatar{"avatar"})

	body := bytes.NewReader(avatar)

	r := httptest.NewRequest("POST", "/api/users/1/avatar", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.SetAvatar(w, r)

	if w.Code != http.StatusCreated {
		t.Error("Status is not 201")
	}
}

func TestFaildSetAvatarNoRights(t *testing.T) {
	t.Parallel()

	h := NewNotEmptyUsersHandler()
	h.Sessions["username"] = 1
	h.Sessions["other username"] = 2
	h.UserAvatars[1] = "avatar"

	avatar, _ := json.Marshal(Avatar{"new avatar"})

	body := bytes.NewReader([]byte(avatar))

	r := httptest.NewRequest("POST", "/api/users/1/avatar", body)
	cookie := &http.Cookie {
		Name: "session_id",
		Value: "other username",
		Expires: time.Now().Add(time.Hour),
	}
	r.AddCookie(cookie)
	r = mux.SetURLVars(r, map[string]string{"user_id": "1"})

	w := httptest.NewRecorder()

	h.SetAvatar(w, r)

	if w.Code != http.StatusForbidden{
		t.Error("Status is not 403")
	}

	if h.UserAvatars[1] != "avatar" {
		t.Error("Wrong changes were saved")
	}
}