package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"server/server/models"
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
	Sessions    map[string]uint
	Users       map[string]*models.User
	UserAvatars map[uint]string
	UserSummary map[uint]models.UserSummary
	Mu          sync.RWMutex
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		Sessions: make(map[string]uint, 10),
		Users:    map[string]*models.User {
			"marat1k": {1, "marat1k", "ABCDE12345", "Marat", "Ishimbaev", "m@m.m", "89032909821"},
		},
		UserAvatars: map[uint]string{},
		UserSummary: map[uint]models.UserSummary{},
		Mu:          sync.RWMutex{},
	}
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users/login")

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	type Response struct {
		ID uint `json:"id"`
	}
	log.Println("Sessions available: ", len(api.Sessions))
	session, err := r.Cookie("session_id")
	if err == nil {
		api.Mu.RLock()
		userId, found := api.Sessions[session.Value]
		api.Mu.RUnlock()
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
	api.Mu.RLock()
	user, ok := api.Users[login]
	api.Mu.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if user.Password != data["password"] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	SID := getSID(64)

	api.Mu.Lock()
	api.Sessions[SID] = user.ID
	api.Mu.Unlock()

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

	log.Println("Sessions available: ", len(api.Sessions))
	session, err := r.Cookie("session_id")
	if err == nil {
		api.Mu.RLock()
		userId, found := api.Sessions[session.Value]
		api.Mu.RUnlock()
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

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	api.Mu.Lock()
	if _, ok := api.Sessions[session.Value]; !ok {
		w.WriteHeader(http.StatusUnauthorized)
		api.Mu.Unlock()
		return
	}

	delete(api.Sessions, session.Value)
	api.Mu.Unlock()

	session.Expires = time.Now().AddDate(0, 0, -1)
	session.Path = "/"
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusCreated)
}

func (api *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users")

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	login, found := data["login"]
	if !found || login == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	api.Mu.RLock()
	if _, ok := api.Users[login]; found && ok {
		w.WriteHeader(http.StatusBadRequest)
		api.Mu.RUnlock()
		return
	}
	api.Mu.RUnlock()

	password, found := data["password"]
	if !found || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	firstName := data["first-name"]
	lastName := data["last-name"]
	email := data["email"]
	phoneNumber := data["phone-number"]

	api.Mu.Lock()
	api.Users[login] = &models.User{uint(len(api.Users) + 1), login, password, firstName, lastName, email, phoneNumber}
	api.Mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}