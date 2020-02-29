package models

type User struct {
	ID uint
	Login string
	Password string

	FirstName string
	LastName string
	Email string
	PhoneNumber string
}
