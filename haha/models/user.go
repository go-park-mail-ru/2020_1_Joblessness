package models

import "time"

type Person struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	Gender string `json:"gender,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
}

type Organization struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Name string `json:"name,omitempty"`
	About string `json:"about,omitempty"`
	Site string `json:"site,omitempty"`
}

type UserLogin struct {
	Login string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Favorite struct {
	ID uint64 `json:"id"`
	Tag string `json:"tag"`
	IsPerson bool `json:"isPerson"`
}

type Favorites []*Favorite

type ResponseRole struct {
	ID uint64 `json:"id"`
	Role string `json:"role"`
}