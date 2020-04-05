package searchRepoPostgres

import (
	"database/sql"
	"github.com/kataras/golog"
	"joblessness/haha/models"
	"strconv"
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

type SearchRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

func (r SearchRepository) SearchPersons(request, since, desc string) (result []*models.Person, err error) {
	switch desc {
	case "true":
		desc = "desc"
	default:
		desc = "asc"
	}

	page, _ := strconv.Atoi(since)

	getPersons := 	`SELECT users.id as userId, p.name, tag, avatar
					FROM users
					JOIN person p on users.person_id = p.id
					WHERE name LIKE '%' || $1 || '%'
					      OR tag LIKE '%' || $1 || '%'
					ORDER BY p.name ` + desc + `, registered 
 					LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(getPersons, request, 9, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]*models.Person, 0)

	for rows.Next() {
		var personDB models.Person
		err := rows.Scan(&personDB.ID, &personDB.FirstName, &personDB.Tag, &personDB.Avatar)
		if err != nil {
			return nil, err
		}

		golog.Error("Name in base: ", personDB.FirstName)
		index := strings.Index(personDB.FirstName, " ")
		if index > -1 {
			personDB.LastName = personDB.FirstName[index+1:]
			personDB.FirstName = personDB.FirstName[:index]
		}
		golog.Error("Name after:",  personDB.FirstName, " ", personDB.LastName)

		result= append(result, &personDB)
	}

	return result, nil
}

func (r SearchRepository) SearchOrganizations(request, since, desc string) (result []*models.Organization, err error) {
	switch desc {
	case "true":
		desc = "desc"
	default:
		desc = "asc"
	}

	page, _ := strconv.Atoi(since)

	getOrgs := 	`SELECT users.id as userId, name, tag, avatar
					FROM users
					JOIN organization o on users.organization_id = o.id
					WHERE name LIKE '%' || $1 || '%'
					      OR tag LIKE '%' || $1 || '%'
					ORDER BY o.name ` + desc + `, registered
					LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(getOrgs, request, 9, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]*models.Organization, 0)

	for rows.Next() {
		var orgDB models.Organization
		err := rows.Scan(&orgDB.ID, &orgDB.Name, &orgDB.Tag, &orgDB.Avatar)
		if err != nil {
			return nil, err
		}

		result= append(result, &orgDB)
	}

	return result, nil
}

func (r SearchRepository) SearchVacancies(request, since, desc string) (result []*models.Vacancy, err error) {
	switch desc {
	case "true":
		desc = "desc"
	default:
		desc = "asc"
	}

	page, _ := strconv.Atoi(since)

	getVacancies := `SELECT users.id, o.name, v.id, v.name, v.keywords, v.salary_from, v.salary_to, v.with_tax
					FROM users
					JOIN organization o on users.organization_id = o.id
					JOIN vacancy v on o.id = v.organization_id
					WHERE v.name LIKE '%' || $1 || '%'
					ORDER BY o.name ` + desc + `, v.name
					LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(getVacancies, request, 9, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]*models.Vacancy, 0)

	for rows.Next() {
		var vacancyDB models.Vacancy
		err := rows.Scan(&vacancyDB.Organization.ID, &vacancyDB.Organization.Name, &vacancyDB.ID,
			&vacancyDB.Name, &vacancyDB.Keywords, &vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax)
		if err != nil {
			return nil, err
		}

		result= append(result, &vacancyDB)
	}

	return result, nil
}