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
	fmt.Println("POST /users/login")

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
		Expires: time.Now().Add(time.Hour),
	}
	http.SetCookie(w, cookie)

	fmt.Println(api.users)
}

func (api *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users/logout")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `No session`, 401)
		return
	}

	if _, ok := api.sessions[session.Value]; !ok {
		http.Error(w, `No session`, 401)
		return
	}

	delete(api.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (api *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users")
}