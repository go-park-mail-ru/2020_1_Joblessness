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
	userSummary map[uint]UserSummary
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		sessions: make(map[string]uint, 10),
		users:    map[string]*User {
			"marat1k": {1, "marat1k", "ABCDE12345", "", "", "", ""},
		},
		userAvatars: map[uint]string{},
		userSummary: map[uint]UserSummary{},
	}
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users/login")
	Cors.PrivateApi(&w, r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	type Response struct {
		Status uint `json:"status"`
		ID uint `json:"id,omitempty"`
	}

	user, ok := api.users[data["login"]]
	if !ok {
		jsonData, _ := json.Marshal(Response{Status:404})
		w.Write(jsonData)
		return
	}

	if user.Password != data["password"] {
		jsonData, _ := json.Marshal(Response{Status:404})
		w.Write(jsonData)
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

	jsonData, _ := json.Marshal(Response{200, user.ID})
	w.Write(jsonData)
}

func (api *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users/logout")
	Cors.PrivateApi(&w, r)

	type Response struct {
		Status uint `json:"status"`
	}

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
		return
	}

	if _, ok := api.sessions[session.Value]; !ok {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
		return
	}

	delete(api.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	jsonData, _ := json.Marshal(Response{200})
	w.Write(jsonData)
}

func (api *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users")
	Cors.PrivateApi(&w, r)

	type Response struct {
		Status uint `json:"status"`
	}

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	login := data["login"]
	if _, ok := api.users[login]; ok {
		jsonData, _ := json.Marshal(Response{400})
		w.Write(jsonData)
		return
	}
	password := data["password"]

	firstName := data["first-name"]
	lastName := data["last-name"]
	email := data["email"]
	phoneNumber := data["phone-number"]

	api.users[login] = &User{uint(len(api.users) + 1), login, password, firstName, lastName, email, phoneNumber}

	jsonData, _ := json.Marshal(Response{200})
	w.Write(jsonData)
}