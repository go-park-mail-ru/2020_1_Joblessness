package authRepoPostgres

import (
	"database/sql"
	"errors"
	"github.com/kataras/golog"
	"joblessness/haha/auth/interfaces"
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
	name := p.FirstName + " " + p.LastName

	return &User{
		ID:             p.ID,
		Login:          p.Login,
		Password:       p.Password,
		Tag:            p.Tag,
		Email:          p.Email,
		Phone:          p.Phone,
		Registered:     p.Registered,
		Avatar:         p.Avatar,
	},
	&Person{
		Name:     name,
		Gender:   p.Gender,
		Birthday: p.Birthday,
	}
}

func toPostgresOrg(o *models.Organization) (*User, *Organization) {
	return &User{
		ID:             o.ID,
		Login:          o.Login,
		Password:       o.Password,
		Tag:            o.Tag,
		Email:          o.Email,
		Phone:          o.Phone,
		Registered:     o.Registered,
		Avatar:         o.Avatar,
	},
	&Organization{
		Name: o.Name,
		Site: o.Site,
	}
}

func toModelPerson(u *User, p *Person) *models.Person {
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
		Registered: u.Registered,
		Avatar:     u.Avatar,
		FirstName:  firstName,
		LastName:   lastName,
		Gender:     p.Gender,
		Birthday:   p.Birthday,
	}
}

func toModelOrganization(u *User, o *Organization) *models.Organization {
	return &models.Organization{
		ID:         u.ID,
		Login:      u.Login,
		Password:   u.Password,
		Tag:        u.Tag,
		Email:      u.Email,
		Phone:      u.Phone,
		Registered: u.Registered,
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
		return authInterfaces.ErrUserAlreadyExists
	}
	return nil
}

func (r UserRepository) CreateUser(user *User) (err error) {
	var personIdSql sql.NullInt64
	var orgIdSql sql.NullInt64
	if user.PersonID != 0 {
		personIdSql.Valid = true
		personIdSql.Int64 = int64(user.PersonID)
	} else if user.OrganizationID != 0 {
		orgIdSql.Valid = true
		orgIdSql.Int64 = int64(user.OrganizationID)
	} else {
		return errors.New("inserted id is 0")
	}

	user.Password, err = salt.HashAndSalt(user.Password)

	insertUser := `INSERT INTO users (login, password, organization_id, person_id, email, phone) 
					VALUES(NULLIF($1, ''), NULLIF($2, ''), NULLIF($3, 0), NULLIF($4, 0), $5, $6)`
	_, err = r.db.Exec(insertUser, user.Login, user.Password, user.OrganizationID, user.PersonID, user.Email, user.Phone)

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

func (r UserRepository) CreateOrganization(org *models.Organization) (err error) {
	dbUser, dbOrg := toPostgresOrg(org)

	err = r.db.QueryRow("INSERT INTO organization (name, site) VALUES($1) RETURNING id", dbOrg.Name, dbOrg.Site).
		Scan(&dbUser.OrganizationID)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(dbUser)

	return err
}

func (r UserRepository) Login(login, password, SID string) (userId uint64, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id, password FROM users WHERE login = $1"
	var hashedPwd string
	rows := r.db.QueryRow(checkUser, login)
	err = rows.Scan(&userId, &hashedPwd)
	if err != nil || !salt.ComparePasswords(hashedPwd, password) {
		golog.Error("DB err - ", err)
		return 0, authInterfaces.ErrWrongLogPas
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

func (r UserRepository) GetPerson(userID uint64) (*models.Person, error) {
	user := User{ID: userID}

	getUser := "SELECT login, person_id, email, phone, avatar, tag FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, userID).
		Scan(&user.Login, &user.PersonID, &user.Email, &user.Phone, &user.Avatar, &user.Tag)
	if err != nil {
		golog.Error(err)
		return nil, authInterfaces.ErrUserNotPerson
	}

	var person Person

	getPerson := "SELECT name, gender, birthday FROM person WHERE id = $1;"
	err = r.db.QueryRow(getPerson, user.PersonID).Scan(&person.Name, &person.Gender, &person.Birthday)
	if err != nil {
		return nil, err
	}

	return toModelPerson(&user, &person), nil
}

func (r UserRepository) changeUser(user *User) error {
	changeUser := `UPDATE users 
					SET password = COALESCE(NULLIF($1, ''), password), 
					    tag = COALESCE(NULLIF($2, ''), tag), 
					    email = COALESCE(NULLIF($3, ''), email),
					    phone = COALESCE(NULLIF($4, ''), phone)
					WHERE id = $5;`
	_, err := r.db.Exec(changeUser, user.Password, user.Tag, user.Email, user.Phone, user.ID)

	return err
}

func (r UserRepository) ChangePerson(p models.Person) error {
	user, dbPerson := toPostgresPerson(&p)

	getUser := "SELECT person_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.PersonID)
	if err != nil {
		return authInterfaces.ErrUserNotPerson
	}

	var birthday sql.NullTime
	if p.Birthday.IsZero() {
		birthday.Valid = false
	} else {
		birthday.Time = p.Birthday
		birthday.Valid = true
	}

	changePerson := `UPDATE person 
					SET name = COALESCE(NULLIF($1, ''), name), 
					    gender = COALESCE(NULLIF($2, ''), gender), 
					    birthday = COALESCE($3, birthday) 
					WHERE id = $4;`
	_, err = r.db.Exec(changePerson, dbPerson.Name, dbPerson.Gender, birthday, user.PersonID)
	if err != nil {
		return err
	}
	err = r.changeUser(user)

	return err
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
		return nil, authInterfaces.ErrUserNotOrg
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
		return authInterfaces.ErrUserNotOrg
	}

	changeOrg := `UPDATE organization 
					SET name = COALESCE(NULLIF($1, ''), name),
					    site = COALESCE(NULLIF($2, ''), site)
					WHERE id = $3;`
	_, err = r.db.Exec(changeOrg, dbOrg.Name, dbOrg.Site, user.OrganizationID)
	if err != nil {
		return err
	}
	err = r.changeUser(user)

	return err
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
			Registered:  time.Time{},
		})
	}

	return result, rows.Err()
}
