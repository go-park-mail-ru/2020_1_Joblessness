package postgres

import (
	"database/sql"
	"errors"
	"joblessness/haha/auth"
	"joblessness/haha/models"
	"time"
)

type User struct {
	ID          uint64
	Login       string
	Password    string
	OrganizationID uint64
	PersonID uint64
	Tag         string
	Email       string
	PhoneNumber string
	Registered  time.Time
}

type Person struct {
	ID uint64
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
		return auth.ErrUserAlreadyExists
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

func (r UserRepository) Login(login, password, SID string) (userId uint64, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id FROM users WHERE login = $1 AND password = $2;"
	err = r.db.QueryRow(checkUser, login, password).Scan(&userId)
	if err != nil {
		return 0, err
	}
	if userId == 0 {
		return 0, auth.ErrWrongLogPas
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

func (r UserRepository) SessionExists(sessionId string) (userId uint64, err error) {
	//TODO session_id - pk

	checkUser := "SELECT user_id, expires FROM session WHERE session_id = $1;"
	var expires time.Time
	err = r.db.QueryRow(checkUser, sessionId).Scan(&userId,  &expires)
	if err != nil {
		return 0, err
	}
	if userId == 0 {
		return 0, auth.ErrWrongSID
	}

	if expires.Before(time.Now()) {
		deleteRow := "DELETE FROM session WHERE session_id = $1;"
		_, err = r.db.Exec(deleteRow, sessionId)
		userId = 0
	}

	return userId, err
}

func (r UserRepository) GetPerson(userID uint64) (models.Person, error) {
	user := User{ID: userID}

	getUser := "SELECT login, password, person, email, phone FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, userID).
		Scan(&user.Login, &user.Password, &user.PersonID, &user.Email, &user.PhoneNumber)
	if err != nil {
		return models.Person{}, err
	}

	if user.PersonID == 0 {
		return models.Person{}, errors.New("this user is not a person")
	}

	var person Person

	getPerson := "SELECT name FROM person WHERE id = $1;"
	err = r.db.QueryRow(getPerson, user.PersonID).Scan(&person.Name)
	if err != nil {
		return models.Person{}, err
	}

	return *toModel(&user, &person), nil
}

func (r UserRepository) ChangePerson(p models.Person) error {
	user := User{ID: p.ID}

	getUser := "SELECT person FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.PersonID)
	if err != nil {
		return err
	}

	if user.PersonID == 0 {
		return errors.New("this user is not a person")
	}

	changePerson := "UPDATE person SET name = $1 WHERE id = $2;"
	_, err = r.db.Exec(changePerson, p.FirstName, user.PersonID)
	if err != nil {
		return err
	}

	return nil
}
