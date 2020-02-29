package models

type User struct {
	ID uint `json:"id,omitempty"`
	Login string `json:"login"`
	Password string `json:"password"`

	FirstName string `json:"first-name"`
	LastName string `json:"last-name"`
	Email string `json:"email"`
	PhoneNumber string `json:"phone-number"`
}

type UserLogin struct {
	Login string `json:"login"`
	Password string `json:"password"`
}
