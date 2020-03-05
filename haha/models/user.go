package models

import (
	"errors"
	"joblessness/haha/utils/database"
	"math/rand"
	"time"
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

type Response struct {
	ID int `json:"id"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetSID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
	err = db.QueryRow("INSERT INTO person (name) VALUES($1) RETURNING id", firstName + " " + lastName).Scan(&personId)
	if err != nil {
		return err
	}

	insertUser := `INSERT INTO users (login, password, person, email, phone) 
					VALUES($1, $2, $3, $4, $5)`
	_, err = db.Exec(insertUser, login, password, personId, email, phone)

	return err
}

func Login(login, password, SID string) (userId int, err error) {
	//TODO user_id, session_id уникальные
	db := database.GetDatabase()
	if db == nil {
		return 0, errors.New("No connection to DB")
	}

	checkUser := "SELECT id FROM users WHERE login = $1 AND password = $2;"
	err = db.QueryRow(checkUser, login, password).Scan(&userId)
	if err != nil {
		return 0, err
	}
	if userId == 0 {
		return 0, nil
	}

	insertSession := `INSERT INTO session (user_id, session_id, expires) 
					VALUES($1, $2, $3)`
	_, err = db.Exec(insertSession, userId, SID, time.Now().Add(time.Hour))

	return userId, err
}

func Logout(sessionId string) (err error) {
	//TODO user_id, session_id уникальные
	db := database.GetDatabase()
	if db == nil {
		return errors.New("No connection to DB")
	}

	deleteRow := "DELETE FROM session WHERE session_id = $1;"
	_, err = db.Exec(deleteRow, sessionId)

	return err
}

func SessionExists(sessionId string) (userId int, err error) {
	//TODO session_id - pk
	db := database.GetDatabase()
	if db == nil {
		return 0, errors.New("No connection to DB")
	}

	checkUser := "SELECT user_id, expires FROM session WHERE session_id = $1;"
	var expires time.Time
	err = db.QueryRow(checkUser, sessionId).Scan(&userId,  &expires)
	if err != nil {
		return 0, err
	}
	if userId == 0 {
		return 0, nil
	}

	if expires.Before(time.Now()) {
		deleteRow := "DELETE FROM session WHERE session_id = $1;"
		_, err = db.Exec(deleteRow, sessionId)
		userId = 0
	}

	return userId, err
}