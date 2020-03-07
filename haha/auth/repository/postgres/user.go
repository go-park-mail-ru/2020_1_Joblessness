package postgres

import (
	"database/sql"
	"errors"
	"joblessness/haha/models"
	"time"
)

type User struct {
	ID          uint
	Login       string
	Password    string
	Tag         string
	Email       string
	PhoneNumber string
	Registered  time.Time
}

type Person struct {
	ID uint
	Name string
	Gender string
	Birthday time.Time
}

func toPostgresPerson(u *models.Person) (*User, *Person) {
	return &User{
		ID: u.ID,
		Login: u.Login,
		Password: u.Password,
		Tag: "",
		Email: u.Email,
		PhoneNumber: u.PhoneNumber,
	},
	&Person{
		Name: u.FirstName + " " + u.LastName,
		Gender: "",
		Birthday: time.Time{},
	}
}

func toModel(u *User, p *Person) *models.Person {
	return &models.Person{
		ID:          u.ID,
		Login:       u.Login,
		Password:    u.Password,
		FirstName:   p.Name,
		LastName:    "",
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
	}
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CreatePerson(user *models.Person) (err error) {
	dbUser, dbPerson := toPostgresPerson(user)

	var columnCount int
	checkUser := "SELECT count(*) FROM users WHERE login = $1"
	err = r.db.QueryRow(checkUser, dbUser.Login).Scan(&columnCount)
	if err != nil {
		return err
	}
	if columnCount != 0 {
		return errors.New("Login taken")
	}

	var personId uint64
	err = r.db.QueryRow("INSERT INTO person (name) VALUES($1) RETURNING id", dbPerson.Name).Scan(&personId)
	if err != nil {
		return err
	}

	insertUser := `INSERT INTO users (login, password, person, email, phone) 
					VALUES($1, $2, $3, $4, $5)`
	_, err = r.db.Exec(insertUser, dbUser.Login, dbUser.Password, personId, dbUser.Email, dbUser.PhoneNumber)

	return err
}

func (r UserRepository) Login(login, password, SID string) (userId int, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id FROM users WHERE login = $1 AND password = $2;"
	err = r.db.QueryRow(checkUser, login, password).Scan(&userId)
	if err != nil {
		return 0, err
	}
	if userId == 0 {
		return 0, errors.New("User wasnt found")
	}

	insertSession := `INSERT INTO session (user_id, session_id, expires) 
					VALUES($1, $2, $3)`
	_, err = r.db.Exec(insertSession, userId, SID, time.Now().Add(time.Hour))

	return userId, err
}

func (r UserRepository) Logout(sessionId string) (err error) {
	//TODO user_id, session_id уникальные

	deleteRow := "DELETE FROM session WHERE session_id = $1;"
	_, err = r.db.Exec(deleteRow, sessionId)

	return err
}

func (r UserRepository) SessionExists(sessionId string) (userId int, err error) {
	//TODO session_id - pk

	checkUser := "SELECT user_id, expires FROM session WHERE session_id = $1;"
	var expires time.Time
	err = r.db.QueryRow(checkUser, sessionId).Scan(&userId,  &expires)
	if err != nil {
		return 0, err
	}
	if userId == 0 {
		return 0, nil
	}

	if expires.Before(time.Now()) {
		deleteRow := "DELETE FROM session WHERE session_id = $1;"
		_, err = r.db.Exec(deleteRow, sessionId)
		userId = 0
	}

	return userId, err
}