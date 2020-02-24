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

	//session, err := r.Cookie("session_id")
	//if err == http.ErrNoCookie {
	//	jsonData, _ := json.Marshal(Response{Status:401})
	//	w.Write(jsonData)
	//	return
	//}
	//_ , found := api.sessions[session.Value]
	//if !found {
	//	jsonData, _ := json.Marshal(Response{Status:401})
	//	w.Write(jsonData)
	//	return
	//}

	var currentUser *User
	userId, _ := strconv.Atoi(mux.Vars(r)["user_id"])

	for _, user := range api.users {
		if (*user).ID == uint(userId) {
			currentUser = user
		}
	}

	if currentUser == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userAvatar, found := api.userAvatars[currentUser.ID]
	if !found {
		userAvatar = "default-avatar.jpg"
	}

	type Response struct {
		User UserInfo `json:"user"`
		Summaries []UserSummary `json:"summaries"`
	}

	jsonData, _ := json.Marshal(Response{
		UserInfo{currentUser.FirstName, currentUser.LastName, "", userAvatar},
		[]UserSummary{},
	})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *AuthHandler) ChangeUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /user/{user_id}")
	Cors.PrivateApi(&w, r)

	session, err := r.Cookie("session_id")
	fmt.Println("session cookie: ", session)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, found := api.sessions[session.Value]
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var currentUser *User

	for _, user := range api.users {
		if (*user).ID == uint(userId) {
			currentUser = user
		}
	}

	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var data map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)
//TODO Проверять есть ли все поля
	(*currentUser).LastName = data["last-name"]
	(*currentUser).FirstName = data["first-name"]
	(*currentUser).Password = data["password"]

	w.WriteHeader(http.StatusNoContent)
}

func (api *AuthHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /users/{user_id}/avatar")
	Cors.PrivateApi(&w, r)

	defer r.Body.Close()

	session, err := r.Cookie("session_id")
	fmt.Println("session cookie: ", session)
	if err == http.ErrNoCookie {
		fmt.Println("No session cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, found := api.sessions[session.Value]
	fmt.Println("session id is ", session.Value)
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if 	reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var currentUser *User

	for _, user := range api.users {
		if (*user).ID == userId {
			currentUser = user
		}
	}

	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Problems not with cookies ")
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

	w.WriteHeader(http.StatusCreated)
}