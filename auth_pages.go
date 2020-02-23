package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	ID uint
	Username string
	Password string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getSID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type AuthHandler struct {
	sessions map[string]uint
	users map[string]*User
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		sessions: make(map[string]uint, 10),
		users:    map[string]*User {
			"marat1k": {1, "marat1k", "password"},
		},
	}
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/login request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	user, ok := api.users[r.FormValue("login")]
	if !ok {
		http.Error(w, `Not found`, 404)
		return
	}

	if user.Password != r.FormValue("password") {
		http.Error(w, `Wrong password`, 400)
		return
	}

	SID := getSID(64)

	api.sessions[SID] = user.ID

	cookie := &http.Cookie {
		Name: "session_id",
		Value: SID,
		Expires: time.Now().Add(time.Minute),
	}
	http.SetCookie(w, cookie)
}