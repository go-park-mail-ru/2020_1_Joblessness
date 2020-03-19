package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"joblessness/haha/auth"
	"joblessness/haha/models"
	"joblessness/haha/utils/salt"
	"strings"
	"time"
)

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
	ID uint64
	Name string
	Gender string
	Birthday time.Time
}

type Organization struct {
	ID uint64
	Name string
	Site string
}

func toPostgresPerson(p *models.Person) (*User, *Person) {
	day, err := time.Parse(time.RFC3339, p.Birthday)
	if err != nil {
		day = time.Time{}
	}

	return &User{
		ID:             p.ID,
		Login:          p.Login,
		Password:       p.Password,
		OrganizationID: 0,
		PersonID:       0,
		Tag:            p.Tag,
		Email:          p.Email,
		Phone:          p.Phone,
		Registered:     time.Time{},
		Avatar:         p.Avatar,
	},
	&Person{
		ID:       0,
		Name:     p.FirstName + " " + p.LastName,
		Gender:   p.Gender,
		Birthday: day,
	}
}

func toPostgresOrg(o *models.Organization) (*User, *Organization) {
	return &User{
		ID:             o.ID,
		Login:          o.Login,
		Password:       o.Password,
		OrganizationID: 0,
		PersonID:       0,
		Tag:            o.Tag,
		Email:          o.Email,
		Phone:          o.Phone,
		Registered:     time.Time{},
		Avatar:         o.Avatar,
	},
	&Organization{
		ID:   0,
		Name: o.Name,
		Site: o.Site,
	}
}

func toModelPerson(u *User, p *Person) *models.Person {
	registered := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n", u.Registered.Year(), u.Registered.Month(),
							  u.Registered.Day(), u.Registered.Hour(), u.Registered.Minute(), u.Registered.Second())

	birthday := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n", u.Registered.Year(), u.Registered.Month(),
							u.Registered.Day(), u.Registered.Hour(), u.Registered.Minute(), u.Registered.Second())

	name := strings.Split(p.Name, " ")
	firstName := name[0]
	var lastName string
	if len(name) > 1 {
		lastName = p.Name[(len(p.Name)) -len(name[0]):]
	}

	return &models.Person{
		ID:         u.ID,
		Login:      u.Login,
		Password:   u.Password,
		Tag:        u.Tag,
		Email:      u.Email,
		Phone:      u.Phone,
		Registered: registered,
		Avatar:     u.Avatar,
		FirstName:  firstName,
		LastName:   lastName,
		Gender:     p.Gender,
		Birthday:   birthday,
	}
}

func toModelOrganization(u *User, o *Organization) *models.Organization {
	registered := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n", u.Registered.Year(), u.Registered.Month(),
							  u.Registered.Day(), u.Registered.Hour(), u.Registered.Minute(), u.Registered.Second())

	return &models.Organization{
		ID:         u.ID,
		Login:      u.Login,
		Password:   u.Password,
		Tag:        u.Tag,
		Email:      u.Email,
		Phone:      u.Phone,
		Registered: registered,
		Avatar:     u.Avatar,
		Name:       o.Name,
		Site:       o.Site,
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

func (r UserRepository) DoesUserExists(login string) (err error) {
	var columnCount int
	checkUser := "SELECT count(*) FROM users WHERE login = $1"
	err = r.db.QueryRow(checkUser, login).Scan(&columnCount)

	if err != nil {
		return err
	}

	if columnCount != 0 {
		return auth.ErrUserAlreadyExists
	}
	return nil
}

func (r UserRepository) CreateUser(login, password, email, phone string, personId, orgId uint64) (err error) {
	var personIdSql sql.NullInt64
	var orgIdSql sql.NullInt64
	if personId != 0 {
		personIdSql.Valid = true
		personIdSql.Int64 = int64(personId)
	} else if orgId != 0 {
		orgIdSql.Valid = true
		orgIdSql.Int64 = int64(orgId)
	} else {
		return errors.New("inserted id is 0")
	}

	password, err = salt.HashAndSalt(password)

	insertUser := `INSERT INTO users (login, password, organization_id, person_id, email, phone) 
					VALUES($1, $2, $3, $4, $5, $6)`
	_, err = r.db.Exec(insertUser, login, password, orgIdSql, personIdSql, email, phone)

	return err
}

func (r UserRepository) SaveAvatarLink(link string, userID uint64) (err error) {
	if link == "" {
		return errors.New("no link to save")
	}

	insertUser := `UPDATE users SET avatar = $1 WHERE id = $2;`
	_, err = r.db.Exec(insertUser, link, userID)

	return err
}

func (r UserRepository) CreatePerson(user *models.Person) (err error) {
	dbUser, dbPerson := toPostgresPerson(user)

	var personId uint64
	err = r.db.QueryRow("INSERT INTO person (name) VALUES($1) RETURNING id", dbPerson.Name).Scan(&personId)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(dbUser.Login, dbUser.Password, dbUser.Email, dbUser.Phone, personId, 0)

	return err
}

func (r UserRepository) CreateOrganization(org *models.Organization) (err error) {
	dbUser, dbOrg := toPostgresOrg(org)

	var orgId uint64
	err = r.db.QueryRow("INSERT INTO organization (name) VALUES($1) RETURNING id", dbOrg.Name).Scan(&orgId)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(dbUser.Login, dbUser.Password, dbUser.Email, dbUser.Phone, 0, orgId)

	return err
}

func (r UserRepository) Login(login, password, SID string) (userId uint64, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id, password FROM users WHERE login = $1"
	var hashedPwd string
	rows := r.db.QueryRow(checkUser, login)
	err = rows.Scan(&userId, &hashedPwd)
	if err != nil || !salt.ComparePasswords(hashedPwd, password) {
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
	//TODO session_id - pk, возвращать тип сессии

	checkUser := "SELECT user_id, expires FROM session WHERE session_id = $1;"
	var expires time.Time
	err = r.db.QueryRow(checkUser, sessionId).Scan(&userId,  &expires)
	if err != nil {
		return 0, auth.ErrWrongSID
	}

	if expires.Before(time.Now()) {
		deleteRow := "DELETE FROM session WHERE session_id = $1;"
		_, err = r.db.Exec(deleteRow, sessionId)
		if err != nil {
			return 0, err
		}
		userId = 0
		return userId, auth.ErrWrongSID
	}

	return userId, err
}

func (r UserRepository) GetPerson(userID uint64) (*models.Person, error) {
	user := User{ID: userID}

	getUser := "SELECT login, password, person_id, email, phone, avatar FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, userID).
		Scan(&user.Login, &user.Password, &user.PersonID, &user.Email, &user.Phone, &user.Avatar)
	if err != nil {
		return nil, auth.ErrUserNotPerson
	}

	var person Person

	getPerson := "SELECT name FROM person WHERE id = $1;"
	err = r.db.QueryRow(getPerson, user.PersonID).Scan(&person.Name)
	if err != nil {
		return nil, err
	}

	return toModelPerson(&user, &person), nil
}

func (r UserRepository) ChangePerson(p models.Person) error {
	user, dbPerson := toPostgresPerson(&p)

	getUser := "SELECT person_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.PersonID)
	if err != nil {
		return auth.ErrUserNotPerson
	}

	changePerson := "UPDATE person SET name = $1 WHERE id = $2;"
	_, err = r.db.Exec(changePerson, dbPerson.Name, user.PersonID)
	if err != nil {
		return err
	}

	return nil
}


func (r UserRepository) GetOrganization(userID uint64) (*models.Organization, error) {
	user := User{ID: userID}

	getUser := "SELECT login, password, organization_id, email, phone, avatar FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, userID).
		Scan(&user.Login, &user.Password, &user.OrganizationID, &user.Email, &user.Phone, &user.Avatar)
	if err != nil {
		return nil, err
	}

	if user.OrganizationID == 0 {
		return nil, auth.ErrUserNotOrg
	}

	var org Organization

	getOrg := "SELECT name FROM organization WHERE id = $1;"
	err = r.db.QueryRow(getOrg, user.OrganizationID).Scan(&org.Name)
	if err != nil {
		return nil, err
	}

	return toModelOrganization(&user, &org), nil
}

func (r UserRepository) ChangeOrganization(o models.Organization) error {
	user, dbOrg := toPostgresOrg(&o)

	getUser := "SELECT organization_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.OrganizationID)
	if err != nil {
		return auth.ErrUserNotOrg
	}

	changePerson := "UPDATE organization SET name = $1 WHERE id = $2;"
	_, err = r.db.Exec(changePerson, dbOrg.Name, user.OrganizationID)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetListOfOrgs(page int) (result []models.Organization, err error) {
	getOrgs := `SELECT users.id as userId, name, site
FROM users, organization
WHERE users.organization_id = organization.id
ORDER BY registered desc
LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(getOrgs, (page - 1)*10, 9)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		userId uint64
		name, site string
	)

	for rows.Next() {
		err := rows.Scan(&userId, &name, &site)
		if err != nil {
			return nil, err
		}

		result= append(result, models.Organization{
			ID:          userId,
			Login:       "",
			Password:    "",
			Name:        name,
			Site:        site,
			Email:       "",
			Phone: "",
			Tag:         "",
			Registered:  "",
		})
	}

	return result, rows.Err()
}
