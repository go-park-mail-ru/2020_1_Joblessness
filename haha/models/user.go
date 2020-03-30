package models

import "time"

type Person struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
}

type Organization struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Name string `json:"name,omitempty"`
	Site string `json:"site,omitempty"`
}

type UserLogin struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type Favorite struct {
	ID uint64 `json:"id"`
	Tag string `json:"tag"`
	IsPerson bool `json:"is_person"`
}

type Favorites []*Favorite

type Response struct {
	ID uint64 `json:"id"`
}