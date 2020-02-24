package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type UserPage struct {
	User interface{} `json:"user,omitempty"`
	Summaries interface{} `json:"summaries"`
}

type UserInfo struct {
	Firstname string `json:"first-name,omitempty"`
	Lastname string `json:"last-name,omitempty"`
	Tag string `json:"tag,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type UserSummary struct {
	Title string `json:"title"`
}

func (api *AuthHandler) GetUserPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /user/{user_id}")
	Cors.PrivateApi(&w, r)

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `No session`, 401)
		return
	}
	_ , found := api.sessions[session.Value]
	if !found {
		http.Error(w, `No session`, 401)
		return
	}

	var currentUser *User
	userId, _ := strconv.Atoi(mux.Vars(r)["user_id"])

	for _, user := range api.users {
		if (*user).ID == uint(userId) {
			currentUser = user
		}
	}

	if currentUser == nil {
		http.Error(w, `No user found`, 404)
		return
	}

	userAvatar, found := api.userAvatars[currentUser.ID]
	if !found {
		userAvatar = "default-avatar.jpg"
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(UserPage{
		User: UserInfo{
			Firstname: currentUser.FirstName,
			Lastname:  currentUser.LastName,
			Tag:       "",
			Avatar:    userAvatar,
		},
		Summaries: []UserSummary{},
	})
}

func (api *AuthHandler) SetUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /user/{user_id}")
	Cors.PrivateApi(&w, r)

	session, err := r.Cookie("session_id")
	fmt.Println("session cookie: ", session.Name)
	if err == http.ErrNoCookie {
		http.Error(w, `No session`, 401)
		return
	}
	userId, found := api.sessions[session.Value]
	if !found {
		http.Error(w, `No session`, 401)
		return
	}
	if 	reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		http.Error(w, `Insufficient rights`, 403)
		return
	}

	var currentUser *User

	for _, user := range api.users {
		if (*user).ID == uint(userId) {
			currentUser = user
		}
	}

	if currentUser == nil {
		http.Error(w, `No user found`, 401)
		return
	}

	var data map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)
//TODO Проверять есть ли все поля
	(*currentUser).LastName = data["last-name"]
	(*currentUser).FirstName = data["first-name"]
	(*currentUser).Password = data["password"]
}

func (api *AuthHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/users/{user_id}/avatar")
	Cors.PrivateApi(&w, r)

	defer r.Body.Close()

	session, err := r.Cookie("session_id")
	fmt.Println("session cookie: ", session)
	if err == http.ErrNoCookie {
		http.Error(w, `No session`, 401)
		return
	}
	userId, found := api.sessions[session.Value]
	if !found {
		http.Error(w, `No session`, 401)
		return
	}
	if 	reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		http.Error(w, `Insufficient rights`, 403)
		return
	}

	var currentUser *User

	for _, user := range api.users {
		if (*user).ID == userId {
			currentUser = user
		}
	}

	if currentUser == nil {
		http.Error(w, `No user found`, 401)
		return
	}

	var data map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)

	avatar, found := data["avatar"]
	if !found {
		avatar = "default-avatar.jpg"
	}

	api.userAvatars[userId] = avatar
}