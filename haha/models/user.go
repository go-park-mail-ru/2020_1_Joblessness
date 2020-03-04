package models

import (
	"errors"
	"joblessness/haha/utils/database"
)

type User struct {
	ID uint `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`

	FirstName string `json:"first-name,omitempty"`
	LastName string `json:"last-name,omitempty"`
	Email string `json:"email,omitempty"`
	PhoneNumber string `json:"phone-number,omitempty"`
}

type UserLogin struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

func CreatePerson(login, password, firstName, lastName, email, phone string) (err error) {
	db := database.GetDatabase()
	if db == nil {
		return errors.New("No connection to DB")
	}

	var columnCount int
	checkUser := "SELECT count(*) FROM users WHERE login = $1"
	err = db.QueryRow(checkUser, login).Scan(&columnCount)
	if err != nil {
		return err
	}
	if columnCount != 0 {
		return errors.New("Login taken")
	}

	var personId uint64
	err = db.QueryRow("INSERT INTO person (name) VALUES($1) RETURNING id", firstName + lastName).Scan(&personId)
	if err != nil {
		return err
	}

	insertUser := `INSERT INTO users (login, password, person, email, phone) 
					VALUES($1, $2, $3, $4, $5)`
	_, err = db.Exec(insertUser, login, password, personId, email, phone)

	return err
}
