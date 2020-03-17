package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kataras/golog"
	"joblessness/haha/auth"
	"joblessness/haha/models"
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

func toPostgresPerson(u *models.Person) (*User, *Person) {
	day, err := time.Parse(time.RFC3339, u.Birthday)
	if err != nil {
		day = time.Time{}
	}

	return &User{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
		Tag: u.Tag,
		Email: u.Email,
		Phone: u.PhoneNumber,
	},

	&Person{
		Name: u.FirstName + " " + u.LastName,
		Gender: u.Gender,
		Birthday: day,
	}
}

func toPostgresOrg(o *models.Organization) (*User, *Organization) {
	return &User{
			ID: o.ID,
			Login: o.Login,
			Password: o.Password,
			Tag: o.Tag,
			Email: o.Email,
			Phone: o.PhoneNumber,
		},

		&Organization{
			Name: o.Name,
			Site: o.Site,
		}
}

func toModelPerson(u *User, p *Person) *models.Person {
	person := &models.Person{
		ID:          u.ID,
		Login:       u.Login,
		Password:    u.Password,
		FirstName:   p.Name,
		Tag: u.Tag,
		Registered: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
			u.Registered.Year(), u.Registered.Month(), u.Registered.Day(),
			u.Registered.Hour(), u.Registered.Minute(), u.Registered.Second()),
		Birthday: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
			u.Registered.Year(), u.Registered.Month(), u.Registered.Day(),
			u.Registered.Hour(), u.Registered.Minute(), u.Registered.Second()),
		Gender: p.Gender,
		Email:       u.Email,
		PhoneNumber: u.Phone,
	}

	name := strings.Split(p.Name, " ")
	person.FirstName = name[0]
	if len(name) > 1 {
		person.LastName = p.Name[(len(name[0])-len(p.Name)):]
	}

	return person
}

func toModelOrganization(u *User, o *Organization) *models.Organization {
	return &models.Organization{
		ID:          u.ID,
		Login:       u.Login,
		Password:    u.Password,
		Tag: u.Tag,
		Registered: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
			u.Registered.Year(), u.Registered.Month(), u.Registered.Day(),
			u.Registered.Hour(), u.Registered.Minute(), u.Registered.Second()),
		Site: o.Site,
		Name: o.Name,
		Email:       u.Email,
		PhoneNumber: u.Phone,
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

	if columnCount != 0 {
		return auth.ErrUserAlreadyExists
	}
	return nil
}

func (r UserRepository) CreateUser(login, password, email, phone string, personId, orgId uint64) (err error) {
	personIdSql := sql.NullInt64{Int64: int64(orgId)}
	orgIdSql := sql.NullInt64{Int64: int64(personId)}
	if personId == 0 {
		orgIdSql.Valid = true
	} else {
		personIdSql.Valid = true
	}

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
	golog.Debug(personId)
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
	err = r.db.QueryRow("INSERT INTO organization (name) VALUES($1) RETURNING id", dbOrg.Name).Scan(orgId)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(dbUser.Login, dbUser.Password, dbUser.Email, dbUser.Phone, 0, orgId)

	return err
}

func (r UserRepository) Login(login, password, SID string) (userId uint64, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id FROM users WHERE login = $1 AND password = $2 AND person_id IS NOT NULL"
	err = r.db.QueryRow(checkUser, login, password).Scan(&userId)
	if err != nil {
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

func (r UserRepository) GetPerson(userID uint64) (*models.Person, error) {
	user := User{ID: userID}

	getUser := "SELECT login, password, person_id, email, phone, avatar FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, userID).
		Scan(&user.Login, &user.Password, &user.PersonID, &user.Email, &user.Phone, &user.Avatar)
	if err != nil {
		return nil, err
	}

	if user.PersonID == 0 {
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
	user := User{ID: p.ID}

	getUser := "SELECT person_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.PersonID)
	if err != nil {
		return err
	}

	if user.PersonID == 0 {
		return auth.ErrUserNotPerson
	}

	changePerson := "UPDATE person SET name = $1 WHERE id = $2;"
	_, err = r.db.Exec(changePerson, p.FirstName, user.PersonID)
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
	user := User{ID: o.ID}

	getUser := "SELECT person_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.PersonID)
	if err != nil {
		return err
	}

	if user.OrganizationID == 0 {
		return auth.ErrUserNotOrg
	}

	changePerson := "UPDATE organization SET name = $1 WHERE id = $2;"
	_, err = r.db.Exec(changePerson, o.Name, user.PersonID)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetListOfOrgs(page int) (result []models.Organization, err error) {
	getOrgs := `SELECT users.id, name, site
FROM users, organization
WHERE users.organization_id = organization.id
ORDER BY registered desc
LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(getOrgs, page*10, 9)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

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
			PhoneNumber: "",
			Tag:         "",
			Registered:  "",
		})
	}

	return result, rows.Err()
}
