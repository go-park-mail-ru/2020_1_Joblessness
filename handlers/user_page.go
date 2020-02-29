package handlers

import (
	_models "../models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
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
	log.Println("GET /user/{user_id}")
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

	var currentUser *_models.User
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
		userAvatar = `data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBxANDw4QDxAQEA8QDxAPDg8QDw8QEw8NFREWFhUSExUYHSggGBolGxUWITEhJSkrLi4uFx8zODMtNygtLisBCgoKBQUFDgUFDisZExkrKysrKysrKysrKysrKysrKysrKysrKysrKysrKysrKysrKysrKysrKysrKysrKysrK//AABEIAOEA4QMBIgACEQEDEQH/xAAcAAEAAgIDAQAAAAAAAAAAAAAABgcFCAEDBAL/xABAEAACAgADBAUIBwcEAwAAAAAAAQIDBAURBgcSITFBUWFxEyJCUoGRocEUI2JysbLRMlNzgpKi4TQ1Q8IVJDP/xAAUAQEAAAAAAAAAAAAAAAAAAAAA/8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAwDAQACEQMRAD8AhYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADkyGUZHicbLhw9Up9stNIrxYGOBZGXbqLpJO++Nf2Yria7jKLdLR14iz+lAVGC08Xuk5fVYnn2TiRLPNhsbgk5Sr8pWvTr873rqAjIOWtOn9OZwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAme7jZP/yF3lrV/wCtTJcS/eWL0PDtA9mxG76WMUb8VxQpf7NfRKzv7kW7gcDXhoRrphGEIrRKK09/ad0IKKSS0SSSS6l1aH2AAAA4aOQBCtrdgKMcpWUpU4jnzitIz+8vmU3mmW24O2VV8HCcep9DXbHuNmSM7b7LwzOh6JK+tOVU+tv1X3MCgAduJolTOVdicZwbjKL6mnzOoAAAAAAAAAAAAAAAAAAAAAAAAD7oqdk4witZSkoxXbJvRL3s2O2byqOBwtNEfQiuN+tN822Uxuzy76TmVGq1jSpXy/lWkf7mvcX0ByAAAAAAAAcM5AFU73tneFxx1a6Wq8SkvS9GfyfsKwNlNoMuWLwuIof/ACVyUdeqenmv36GttkHGUovk4txa709APkAAAAAAAAAAAAAAAAAAAAAAOQLL3J0J242z1a64L+aUn/1LZKs3Iy/166/qH+ctMAAAAAAAAAAABrxt1g/IZljILkna5x+7NKXzNhyi97EUs0s066qW/Hh0+QEOAAAAAAAAAAAAAAAAAAAAAAABYO5nFcGNvqb/APrRqu+UJJ/g2XIa4bKZj9Dx2Fu6o2JT+5LzZfB/A2NT10a6H0AfQAAAAAAAAAAFEb1LePNL/swqh7VBP5l7SenN9C5vwNbdpMd9JxmKu6p3Tcfu66L4IDGgAAAAAAAAAAAAAAAAAAAAAAA5L23a5/8ATsFGMnrdh9KrF1uOnmy9xRBl9l8+sy3Ewurfm/s2w6rK9ea+YGxwPHleYV4uqF1UlKE1qu59j70ewAAAAAABg677o1wlObUYxTlJvoSQEd3hZx9CwFsk9LLfqalrz4pdLXgtWUASXbvaV5niW4tqivWNEe7rm+9kaAAAAAAAAAAAAAAAAAAAAAAAB24fDztko1wlOT6Ixi2wOoErwG7zMbo8Xko1680rJqL08FqeLPdkMbgFx3Va19dlb4kvHsA7dkdrr8snpHz6JPz6m+XjHsZduQZ9h8wrVlE03p50PSi+xo1vPTgcdbhpqymyVc10OL0A2cBTuUb1MRWlHEVxt05cUfNkySYXengpL6yFsH18lICfAg9u9DAJeb5WT7ODQwWZ72G01hqNH61j1+CAs3H42vD1ystnGEIrVtvT3FNbdbdyx+tGH1hh/SfRK3n19iI5ne0GKx8uLEWuaXRBcox8EYoAD3ZTlN+Ns8nh65Tl0vToiu1vqJLbuzzGK14an3Kzn+AEMB7syyfE4R6X02V9Wri+H+pcjwgAAAAAAAAAAAAAAAAD1ZdgLcVYq6YSnN9CS+LfUj6ynLrMZdXTUtZzlovsrrZfmymzFOWUqFaTsa1ttf7U5fJAQfId1TfDPG297qrfwcv0LEynJMNg4qGHphWu1Javvb6WZA5A40Pi2pTi4ySlGSakmtU12NHYAKk2x3bTrcrsCnOHTKj0o9fmdvgVvZBxbjJOMk9GmtGn2M2iI3tJsVhMx1lOHk7tOV1fKX8y6Je0DX8E6zbdfjaW3RKF8OemnmS074vl8SNYvZzG0vSeFvXhXKS961AxQPfVkuKm9I4a9v8Ag2foZ3Ld3uY3ta1eRj61skvgtQImSfZPYvEZlJS0dWH5a3SWnF9xdf4FgbObscPhmp4mX0ma0ajpw1xfbp0v2k7rrUUlFJJLRJLRJdwGOyDIqMvpVVENF0yk+cpy7ZPrZkzkAdV+HhbFxnGM4taNSSaa9pB8+3YYW/inh28PN8+FedW393q9hPQwNddotl8Vl0vroaw9G2Org/0MIbPYzCwvhKuyEZwktJRktU0UVt5sq8su1hzw9jbrfqv1GBFgcnAAAAAAAAAAA5S16OnoXiBaW5vJ+V2Lkul+Tq9n7TX4FpGG2Ry9YXA4arTRqtSl3yfNszIAAAAAAAAHA0OQBxoDkAAAAAAAAADBbZZMsfg7qtPPUXOt9liWvIzoYGrk4uLafJp6NdjXJnyZ7bnAfRswxMEtIubnHwlzMCAAAAAAAAAMjs7hfL4vDV+tbHXwT1McS7ddh/KZlU2tVCMpfAC9YR0SXYkj6AAAAAAAAAAAAAAAAAAAAAAAAAApzfLhOHF02fvK9H4plelvb58NxYfD2ac42OLfc0VCAAAAAAAAALE3MUa4rET9WpJe1ldlsblafq8XZ2zhH4agWaAAAAAAAAAAAAAAAAAAAAAAAAAAIjvRw/lMsufXCUJ+GjKINjtrqPK4DFx7aZv3LX5GuIAAAAAAAAAuvc/Rw5fKXr3S+GiKUL83aUeTyvDfb4p++X+AJSAAAAAAAAAAAAAAAAAAAAAAAAAAOjHV8dVsfWrnH3xZrLfDhnNdkpL3Nm0LNatoKPJYvEw9W6a/u/yBjwAAAAAAAGbGbH1cGX4OPZTH48zXRdKNl8mhw4bDrspr/KgPaAAAAAAAAAAAAAAAAAAAAAAAAAABr1t7VwZnjF22uXvNhSht6FfDml/eoS98QImAAAAAAAD6j0rxRszlf+nw/wDBr/IgAPUAAAAAAAAAAAAAAAAAAAAAAAAAABRO9T/dLvuV/lQAEQAAAAAf/9k=`
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
	log.Println("PUT /user/{user_id}")
	Cors.PrivateApi(&w, r)

	session, err := r.Cookie("session_id")
	log.Println("session cookie: ", session)
	if err != nil {
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

	var currentUser *_models.User

	log.Println("Users counter", len(api.users))
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
	//(*currentUser).Password = data["password"]

	w.WriteHeader(http.StatusNoContent)
}

func (api *AuthHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users/{user_id}/avatar")
	Cors.PrivateApi(&w, r)

	defer r.Body.Close()

	session, err := r.Cookie("session_id")
	log.Println("session cookie: ", session)
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

	var currentUser *_models.User

	for _, user := range api.users {
		if (*user).ID == userId {
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

	avatar, found := data["avatar"]
	if !found {
		avatar = `data:image/jpeg;base64,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`
	}

	api.userAvatars[userId] = avatar

	w.WriteHeader(http.StatusCreated)
}