package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	ID uint
	Login string
	Password string

	FirstName string
	LastName string
	Email string
	PhoneNumber string
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
	userAvatars map[uint]string
	userSummary map[uint]Summary
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		sessions: make(map[string]uint, 10),
		users:    map[string]*User {
			"marat1k": {1, "marat1k", "password", "", "", "", ""},
		},
		userAvatars: map[uint]string{},
		userSummary: map[uint]Summary{},
	}
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users/login")

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	user, ok := api.users[data["login"]]
	if !ok {
		http.Error(w, `Not found`, 404)
		return
	}

	if user.Password != data["password"] {
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

	type Response struct {
		ID uint `json:"id"`
	}
	jsonData, _ := json.Marshal(Response{user.ID})
	w.Write(jsonData)
}

func (api *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users/logout")

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

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	login := data["login"]
	if _, ok := api.users[login]; ok {
		http.Error(w, `Login already exist`, 400)
		return
	}
	password := data["password"]

	firstName := data["first-name"]
	lastName := data["last-name"]
	email := data["email"]
	phoneNumber := data["phone-number"]

	api.users[login] = &User{uint(len(api.users) + 1), login, password, firstName, lastName, email, phoneNumber}

	fmt.Println(api.users)
}