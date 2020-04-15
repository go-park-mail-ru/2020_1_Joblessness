package authPostgres

import (
	"database/sql"
	"fmt"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/models"
	"joblessness/haha/utils/salt"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

type User struct {
	ID             uint64
	Login          string
	Password       string
	OrganizationID uint64
	PersonID       uint64
	Tag            string
	Email          string
	Phone          string
	Registered     time.Time
	Avatar         string
}

type Person struct {
	ID       uint64
	Name     string
	Gender   string
	Birthday time.Time
}

type Organization struct {
	ID    uint64
	Name  string
	Site  string
	About string
}

func toPostgresPerson(p *models.Person) (*User, *Person) {
	name := p.FirstName
	if p.LastName != "" {
		name += " " + p.LastName
	}

	return &User{
			ID:         p.ID,
			Login:      p.Login,
			Password:   p.Password,
			Tag:        p.Tag,
			Email:      p.Email,
			Phone:      p.Phone,
			Registered: p.Registered,
			Avatar:     p.Avatar,
		},
		&Person{
			Name:     name,
			Gender:   p.Gender,
			Birthday: p.Birthday,
		}
}

func toPostgresOrg(o *models.Organization) (*User, *Organization) {
	return &User{
			ID:         o.ID,
			Login:      o.Login,
			Password:   o.Password,
			Tag:        o.Tag,
			Email:      o.Email,
			Phone:      o.Phone,
			Registered: o.Registered,
			Avatar:     o.Avatar,
		},
		&Organization{
			Name:  o.Name,
			Site:  o.Site,
			About: o.About,
		}
}

func (r AuthRepository) CreateUser(user *User) (err error) {
	user.Password, err = salt.HashAndSalt(user.Password)

	insertUser := `INSERT INTO users (login, password, organization_id, person_id, email, phone, tag) 
					VALUES(NULLIF($1, ''), NULLIF($2, ''), NULLIF($3, 0), NULLIF($4, 0), $5, $6, $7)`
	_, err = r.db.Exec(insertUser, user.Login, user.Password, user.OrganizationID,
		user.PersonID, user.Email, user.Phone, user.Tag)

	return err
}

func (r *AuthRepository) CreatePerson(user *models.Person) (err error) {
	dbUser, dbPerson := toPostgresPerson(user)
	if dbPerson.Birthday.IsZero() {
		dbPerson.Birthday.AddDate(1950, 0, 0)
	}

	err = r.db.QueryRow("INSERT INTO person (name, gender, birthday) VALUES($1, $2, $3) RETURNING id",
		dbPerson.Name, dbPerson.Gender, dbPerson.Birthday).
		Scan(&dbUser.PersonID)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(dbUser)

	return err
}

func (r *AuthRepository) CreateOrganization(org *models.Organization) (err error) {
	dbUser, dbOrg := toPostgresOrg(org)

	err = r.db.QueryRow("INSERT INTO organization (name, site, about) VALUES($1, $2, $3) RETURNING id",
		dbOrg.Name, dbOrg.Site, dbOrg.About).
		Scan(&dbUser.OrganizationID)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(dbUser)

	return err
}

func (r *AuthRepository) Login(login, password, SID string) (userId uint64, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id, password FROM users WHERE login = $1"
	var hashedPwd string
	rows := r.db.QueryRow(checkUser, login)
	err = rows.Scan(&userId, &hashedPwd)
	if err != nil || !salt.ComparePasswords(hashedPwd, password) {
		return 0, fmt.Errorf("%w, login: %s, password: %s", authInterfaces.ErrWrongLoginOrPassword, login, password)
	}

	insertSession := `INSERT INTO session (user_id, session_id, expires) 
					VALUES($1, $2, $3)`
	_, err = r.db.Exec(insertSession, userId, SID, time.Now().Add(time.Hour))

	return userId, err
}

func (r *AuthRepository) Logout(sessionId string) (err error) {
	//TODO user_id, session_id уникальные

	deleteRow := "DELETE FROM session WHERE session_id = $1;"
	_, err = r.db.Exec(deleteRow, sessionId)

	return err
}

func (r *AuthRepository) SessionExists(sessionId string) (userId uint64, err error) {
	//TODO session_id - pk, возвращать тип сессии

	checkUser := "SELECT user_id, expires FROM session WHERE session_id = $1;"
	var expires time.Time
	err = r.db.QueryRow(checkUser, sessionId).Scan(&userId, &expires)
	if err != nil {
		return 0, authInterfaces.ErrWrongSID
	}

	if expires.Before(time.Now()) {
		deleteRow := "DELETE FROM session WHERE session_id = $1;"
		_, err = r.db.Exec(deleteRow, sessionId)
		if err != nil {
			return 0, err
		}
		userId = 0
		return userId, authInterfaces.ErrWrongSID
	}

	return userId, err
}

func (r *AuthRepository) DoesUserExists(login string) (err error) {
	var columnCount int
	checkUser := "SELECT count(*) FROM users WHERE login = $1"
	err = r.db.QueryRow(checkUser, login).Scan(&columnCount)

	if err != nil {
		return err
	}

	if columnCount != 0 {
		return fmt.Errorf("%w, login: %s", authInterfaces.ErrUserAlreadyExists, login)
	}
	return nil
}

func (r *AuthRepository) GetRole(userID uint64) (string, error) {
	var personID, organizationID sql.NullInt64
	checkUser := "SELECT person_id, organization_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(checkUser, userID).Scan(&personID, &organizationID)
	if err != nil {
		return "", err
	}

	if personID.Valid {
		return "person", nil
	} else if organizationID.Valid {
		return "organization", nil
	}
	return "", nil
}
