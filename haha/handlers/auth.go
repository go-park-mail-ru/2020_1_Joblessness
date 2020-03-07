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
	Users       map[string]*models.Person
	UserAvatars map[uint]string
	UserSummary map[uint]models.UserSummary
	Mu          sync.RWMutex
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		Sessions: make(map[string]uint, 10),
		Users:    map[string]*models.Person{
			"marat1k": {1, "marat1k", "ABCDE12345", "Marat", "Ishimbaev", "m@m.m", "89032909821"},
		},
		UserAvatars: map[uint]string{},
		UserSummary: map[uint]models.UserSummary{},
		Mu:          sync.RWMutex{},
	}
}



func (api *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.Person
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