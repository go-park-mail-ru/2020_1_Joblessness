package models

type Person struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Avatar string `json:"avatar,omitempty"`

	FirstName string `json:"first-name,omitempty"`
	LastName string `json:"last-name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Email string `json:"email,omitempty"`
	PhoneNumber string `json:"phone-number,omitempty"`
	Tag string `json:"tag,omitempty"`
	Birthday string `json:"birth_date,omitempty"`
	Registered string `json:"registered,omitempty"`
}

type Organization struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Avatar string `json:"avatar,omitempty"`

	Name string `json:"name,omitempty"`
	Site string `json:"site,omitempty"`
	Email string `json:"email,omitempty"`
	PhoneNumber string `json:"phone-number,omitempty"`
	Tag string `json:"tag,omitempty"`
	Registered string `json:"registered,omitempty"`
}

type UserLogin struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	ID uint64 `json:"id"`
}