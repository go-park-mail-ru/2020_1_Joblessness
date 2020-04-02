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
	About string
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
		About: o.About,
	}
}

func toModelPerson(u *User, p *Person) *models.Person {
	var lastName, firstName string
	index := strings.Index(p.Name, " ")
	if index > -1 {
		lastName = p.Name[index+1:]
		firstName = p.Name[:index]
	} else {
		firstName = p.Name
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
		About:      o.About,
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

	insertUser := `INSERT INTO users (login, password, organization_id, person_id, email, phone, tag) 
					VALUES(NULLIF($1, ''), NULLIF($2, ''), NULLIF($3, 0), NULLIF($4, 0), $5, $6, $7)`
	_, err = r.db.Exec(insertUser, user.Login, user.Password, user.OrganizationID,
		user.PersonID, user.Email, user.Phone, user.Tag)

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

func (r UserRepository) Login(login, password, SID string) (userId uint64, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id, password FROM users WHERE login = $1"
	var hashedPwd string
	rows := r.db.QueryRow(checkUser, login)
	err = rows.Scan(&userId, &hashedPwd)
	if err != nil || !salt.ComparePasswords(hashedPwd, password) {
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

func (r UserRepository) GetRole(userID uint64) (string, error) {
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

	getOrg := "SELECT name, site, about FROM organization WHERE id = $1;"
	err = r.db.QueryRow(getOrg, user.OrganizationID).Scan(&org.Name, &org.Site, &org.About)
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
					    site = COALESCE(NULLIF($2, ''), site),
					    about = COALESCE(NULLIF($3, ''), about)
					WHERE id = $4;`
	_, err = r.db.Exec(changeOrg, dbOrg.Name, dbOrg.Site, dbOrg.About, user.OrganizationID)
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

	rows, err := r.db.Query(getOrgs, page*10, 9)

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

func (r UserRepository) SetOrDeleteLike(userID, favoriteID uint64) (res bool, err error) {
	setLike := `INSERT INTO favorite (user_id, favorite_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING;`
	rows, err := r.db.Exec(setLike, userID, favoriteID)
	if err != nil {
		return false, err
	}
	if rowsAf, err := rows.RowsAffected(); rowsAf != 0 {
		return true, err
	}

	deleteLike := `DELETE FROM favorite 
				WHERE user_id = $1 AND favorite_id = $2;`
	_, err = r.db.Exec(deleteLike, userID, favoriteID)

	return false, err
}

func (r UserRepository) GetUserFavorite(userID uint64) (res models.Favorites, err error) {
	getFavorite := `SELECT u.id, u.tag, u.person_id
				FROM favorite f, users u 
				WHERE f.favorite_id = u.id
				AND f.user_id = $1;`
	rows, err := r.db.Query(getFavorite, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		personID sql.NullInt64
	)
	for rows.Next() {
		var favorite models.Favorite
		err := rows.Scan(&favorite.ID, &favorite.Tag, &personID)
		if err != nil {
			return nil, err
		}

		if personID.Valid {
			favorite.IsPerson = true
		} else {
			favorite.IsPerson = false
		}
		res = append(res, &favorite)
	}

	return res, err
}
