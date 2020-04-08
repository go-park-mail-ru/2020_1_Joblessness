package models

import "time"

type Person struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty" validate:"required,max=60"`
	Password string `json:"password,omitempty" validate:"required,max=60"`
	Tag string `json:"tag,omitempty" validate:"omitempty,min=6,max=20"`
	Email string `json:"email,omitempty" validate:"omitempty,max=60,email"`
	Phone string `json:"phone,omitempty" validate:"max=20"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	FirstName string `json:"firstName,omitempty" validate:"max=60"`
	LastName string `json:"lastName,omitempty" validate:"max=60"`
	Gender string `json:"gender,omitempty" validate:"max=10"`
	Birthday time.Time `json:"birthday,omitempty"`
}

type Organization struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty" validate:"required,max=60"`
	Password string `json:"password,omitempty" validate:"required,max=60"`
	Tag string `json:"tag,omitempty" validate:"omitempty,min=6,max=20"`
	Email string `json:"email,omitempty" validate:"max=60"`
	Phone string `json:"phone,omitempty" validate:"max=20"`
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

type Role struct {
	Person bool
	Organization bool
}
