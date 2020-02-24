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

	type Data struct {
		User UserInfo `json:"user"`
		Summaries []UserSummary `json:"summaries"`
	}

	type Response struct {
		Status uint `json:"status"`
		Data Data `json:"data,omitempty"`
	}

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		jsonData, _ := json.Marshal(Response{Status:401})
		w.Write(jsonData)
		return
	}
	_ , found := api.sessions[session.Value]
	if !found {
		jsonData, _ := json.Marshal(Response{Status:401})
		w.Write(jsonData)
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
		jsonData, _ := json.Marshal(Response{Status:401})
		w.Write(jsonData)
		return
	}

	userAvatar, found := api.userAvatars[currentUser.ID]
	if !found {
		userAvatar = "default-avatar.jpg"
	}

	jsonData, _ := json.Marshal(Response{200, Data{
		UserInfo{currentUser.FirstName, currentUser.LastName, "", userAvatar},
		[]UserSummary{},
	}})
	w.Write(jsonData)
}

func (api *AuthHandler) SetUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /user/{user_id}")
	Cors.PrivateApi(&w, r)

	type Response struct {
		Status uint `json:"status"`
	}

	session, err := r.Cookie("session_id")
	fmt.Println("session cookie: ", session)
	if err == http.ErrNoCookie {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
		return
	}
	userId, found := api.sessions[session.Value]
	if !found {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
		return
	}
	if reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		jsonData, _ := json.Marshal(Response{403})
		w.Write(jsonData)
		return
	}

	var currentUser *User

	for _, user := range api.users {
		if (*user).ID == uint(userId) {
			currentUser = user
		}
	}

	if currentUser == nil {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
		return
	}

	var data map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)
//TODO Проверять есть ли все поля
	(*currentUser).LastName = data["last-name"]
	(*currentUser).FirstName = data["first-name"]
	(*currentUser).Password = data["password"]

	jsonData, _ := json.Marshal(Response{200})
	w.Write(jsonData)
}

func (api *AuthHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/users/{user_id}/avatar")
	Cors.PrivateApi(&w, r)

	defer r.Body.Close()

	type Response struct {
		Status uint `json:"status"`
	}

	session, err := r.Cookie("session_id")
	fmt.Println("session cookie: ", session)
	if err == http.ErrNoCookie {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
		return
	}
	userId, found := api.sessions[session.Value]
	if !found {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
		return
	}
	if 	reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		jsonData, _ := json.Marshal(Response{403})
		w.Write(jsonData)
		return
	}

	var currentUser *User

	for _, user := range api.users {
		if (*user).ID == userId {
			currentUser = user
		}
	}

	if currentUser == nil {
		jsonData, _ := json.Marshal(Response{401})
		w.Write(jsonData)
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

	jsonData, _ := json.Marshal(Response{200})
	w.Write(jsonData)
}