package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserPage struct {
	User interface{} `json:"user,omitempty"`
	Summaries interface{} `json:"summaries"`
}

type UserInfo struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname string `json:"lastname,omitempty"`
	Tag string `json:"tag,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type UserAvatar struct {
	Avatar string `json:"avatar,omitempty"`
}

type Summary struct {
	Title string `json:"title"`
}

func (api *AuthHandler) configureSummary(id uint) (userSummaries []Summary) {
	userSummary, found := api.userSummary[id]
	if !found {
		return []Summary{}
	}

	userSummaries = append(userSummaries, userSummary)

	return userSummaries
}

func (api *AuthHandler) GetUserPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /user/{id}")
	Cors.EnableCors(&w, r)

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `No session`, 401)
		return
	}
	userId, found := api.sessions[session.Value]
	if !found {
		http.Error(w, `No session`, 401)
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
		Summaries: api.configureSummary(userId),
	})
}

func (api *AuthHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/users/{id}/avatar")
	Cors.EnableCors(&w, r)

	defer r.Body.Close()

	session, err := r.Cookie("session_id")
	fmt.Println("session cook", session)
	if err == http.ErrNoCookie {
		http.Error(w, `No session`, 401)
		return
	}
	userId, found := api.sessions[session.Value]
	if !found {
		http.Error(w, `No session`, 401)
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

	decoder := json.NewDecoder(r.Body)
	userAvatar := new(UserAvatar)

	err = decoder.Decode(userAvatar)
	if err != nil {
		fmt.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	api.userAvatars[userId] = (*userAvatar).Avatar
}