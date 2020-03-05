package handlers

import (
	"encoding/json"
	"io/ioutil"
	"joblessness/haha/models"
	"log"
	"net/http"
	"sync"
	"time"
)

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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.UserLogin
	err = json.Unmarshal(body, &user)
	log.Println("user recieved: ", user)
	if err != nil {
		log.Println("Unmarshal went wrong")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		log.Println("login or password empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	SID := models.GetSID(64)
	userId, err := models.Login(user.Login, user.Password, SID)
	if err != nil {
		log.Println("db broken ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

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

	jsonData, _ := json.Marshal(models.Response{userId})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (api *AuthHandler) Check(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users/check")

	log.Println("Sessions available: ", len(api.Sessions))
	_, err := r.Cookie("session_id")
	if err == nil {
		userId := 1
		if true {
			jsonData, _ := json.Marshal(models.Response{userId})
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	log.Println("user recieved: ", user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = models.CreatePerson(user.Login, user.Password, user.FirstName, user.LastName, user.Email, user.PhoneNumber)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}