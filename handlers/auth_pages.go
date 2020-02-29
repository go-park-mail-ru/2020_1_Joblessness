package handlers

import (
	_models "../models"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

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
	users map[string]*_models.User
	userAvatars map[uint]string
	userSummary map[uint]UserSummary
	mu sync.RWMutex
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		sessions: make(map[string]uint, 10),
		users:    map[string]*_models.User {
			"marat1k": {1, "marat1k", "ABCDE12345", "Marat", "Ishimbaev", "m@m.m", "89032909821"},
		},
		userAvatars: map[uint]string{},
		userSummary: map[uint]UserSummary{},
		mu: sync.RWMutex{},
	}
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users/login")
	Cors.PrivateApi(&w, r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	type Response struct {
		ID uint `json:"id"`
	}
	log.Println("Sessions available: ", len(api.sessions))
	session, err := r.Cookie("session_id")
	if err == nil {
		api.mu.RLock()
		userId, found := api.sessions[session.Value]
		api.mu.RUnlock()
		if found {
			jsonData, _ := json.Marshal(Response{userId})
			w.WriteHeader(http.StatusCreated)
			w.Write(jsonData)
			return
		}
	}

	login, found := data["login"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	api.mu.RLock()
	user, ok := api.users[login]
	api.mu.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if user.Password != data["password"] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	SID := getSID(64)

	api.mu.Lock()
	api.sessions[SID] = user.ID
	api.mu.Unlock()

	cookie := &http.Cookie {
		Name: "session_id",
		Value: SID,
		Expires: time.Now().Add(time.Hour),
		MaxAge: 100000,
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	jsonData, _ := json.Marshal(Response{user.ID})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (api *AuthHandler) Check(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users/check")
	Cors.PrivateApi(&w, r)

	log.Println("Sessions available: ", len(api.sessions))
	session, err := r.Cookie("session_id")
	if err == nil {
		api.mu.RLock()
		userId, found := api.sessions[session.Value]
		api.mu.RUnlock()
		if found {
			type Response struct {
				ID uint `json:"id"`
			}
			jsonData, _ := json.Marshal(Response{userId})
			w.WriteHeader(http.StatusCreated)
			w.Write(jsonData)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (api *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users/logout")
	Cors.PrivateApi(&w, r)

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	api.mu.Lock()
	if _, ok := api.sessions[session.Value]; !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	delete(api.sessions, session.Value)
	api.mu.Unlock()

	session.Expires = time.Now().AddDate(0, 0, -1)
	session.Path = "/"
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusCreated)
}

func (api *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users")
	Cors.PrivateApi(&w, r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	login, found := data["login"]
	if !found || login == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	api.mu.RLock()
	if _, ok := api.users[login]; found && ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	api.mu.RUnlock()

	password, found := data["password"]
	if !found || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	firstName := data["first-name"]
	lastName := data["last-name"]
	email := data["email"]
	phoneNumber := data["phone-number"]

	api.mu.Lock()
	api.users[login] = &_models.User{uint(len(api.users) + 1), login, password, firstName, lastName, email, phoneNumber}
	api.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}