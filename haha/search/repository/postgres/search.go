package searchPostgres

import (
	"database/sql"
	"joblessness/haha/models/base"
	"strconv"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

func (r SearchRepository) SearchPersons(request, since, desc string) (result []*baseModels.Person, err error) {
	switch desc {
	case "true":
		desc = "desc"
	default:
		desc = "asc"
	}

	page, _ := strconv.Atoi(since)

	getPersons := `SELECT users.id as userId, p.name, p.surname, tag, avatar
					FROM users
					JOIN person p on users.person_id = p.id
					WHERE lower(name) LIKE lower('%' || $1 || '%')
					      OR lower(tag) LIKE lower('%' || $1 || '%')
					ORDER BY p.name ` + desc + `, registered 
 					LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(getPersons, request, 10, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]*baseModels.Person, 0)

	for rows.Next() {
		var personDB baseModels.Person
		err := rows.Scan(&personDB.ID, &personDB.FirstName, &personDB.LastName, &personDB.Tag, &personDB.Avatar)
		if err != nil {
			return nil, err
		}

		result = append(result, &personDB)
	}

	return result, nil
}

func (r SearchRepository) SearchOrganizations(request, since, desc string) (result []*baseModels.Organization, err error) {
	switch desc {
	case "true":
		desc = "desc"
	default:
		desc = "asc"
	}

	page, _ := strconv.Atoi(since)

	getOrgs := `SELECT users.id as userId, name, tag, avatar
					FROM users
					JOIN organization o on users.organization_id = o.id
					WHERE lower(name) LIKE lower('%' || $1 || '%')
					      OR lower(tag) LIKE lower('%' || $1 || '%')
					ORDER BY o.name ` + desc + `, registered
					LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(getOrgs, request, 10, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]*baseModels.Organization, 0)

	for rows.Next() {
		var orgDB baseModels.Organization
		err := rows.Scan(&orgDB.ID, &orgDB.Name, &orgDB.Tag, &orgDB.Avatar)
		if err != nil {
			return nil, err
		}

		result = append(result, &orgDB)
	}

	return result, nil
}

func (r SearchRepository) SearchVacancies(request, since, desc string) (result []*baseModels.Vacancy, err error) {
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
					JOIN vacancy v on users.id = v.organization_id
					WHERE lower(v.name) LIKE lower('%' || $1 || '%')
					ORDER BY o.name ` + desc + `, v.name
					LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(getVacancies, request, 10, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]*baseModels.Vacancy, 0)

	for rows.Next() {
		var vacancyDB baseModels.Vacancy
		err := rows.Scan(&vacancyDB.Organization.ID, &vacancyDB.Organization.Name, &vacancyDB.ID,
			&vacancyDB.Name, &vacancyDB.Keywords, &vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax)
		if err != nil {
			return nil, err
		}

		result = append(result, &vacancyDB)
	}

	return result, nil
}
