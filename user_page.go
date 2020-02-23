package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
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

type UserPageHandler struct {
	usersInfo map[string]UserAvatar
}

func (handler *UserPageHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /user/{id}")
	Cors.EnableCors(&w, r)

}