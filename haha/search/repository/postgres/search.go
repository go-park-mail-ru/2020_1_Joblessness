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

func (r SearchRepository) SearchPersons(params *baseModels.SearchParams) (result []*baseModels.Person, err error) {
	switch params.Desc {
	case "true":
		params.Desc = "desc"
	default:
		params.Desc = "asc"
	}

	page, _ := strconv.Atoi(params.Since)

	getPersons := `SELECT users.id as userId, p.name, p.surname, tag, avatar
					FROM users
					JOIN person p on users.person_id = p.id
					WHERE to_tsvector('russian', p.name) @@ plainto_tsquery('russian', $1)
						  OR to_tsvector('russian', p.surname) @@ plainto_tsquery('russian', $1)
					      OR lower(tag) LIKE lower('%' || $1 || '%')
						  OR $1 = ''
					ORDER BY p.name ` + params.Desc + `, registered 
 					LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(getPersons, params.Request, 10, page*10)
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

func (r SearchRepository) SearchOrganizations(params *baseModels.SearchParams) (result []*baseModels.Organization, err error) {
	switch params.Desc {
	case "true":
		params.Desc = "desc"
	default:
		params.Desc = "asc"
	}

	page, _ := strconv.Atoi(params.Since)

	getOrgs := `SELECT users.id as userId, name, tag, avatar
					FROM users
					JOIN organization o on users.organization_id = o.id
					WHERE to_tsvector('russian', o.name) @@ plainto_tsquery('russian', $1)
					      OR lower(tag) LIKE lower('%' || $1 || '%')
						  OR $1 = ''
					ORDER BY o.name ` + params.Desc + `, registered
					LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(getOrgs, params.Request, 10, page*10)
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

func (r SearchRepository) SearchVacancies(params *baseModels.SearchParams) (result []*baseModels.Vacancy, err error) {
	switch params.Desc {
	case "true":
		params.Desc = "desc"
	default:
		params.Desc = "asc"
	}

	page, _ := strconv.Atoi(params.Since)

	getVacancies := `SELECT users.id, users.avatar, o.name, v.id, v.name, v.keywords, v.salary_from, v.salary_to, v.with_tax
					FROM users
					JOIN organization o on users.organization_id = o.id
					JOIN vacancy v on users.id = v.organization_id
					WHERE to_tsvector('russian', v.name) @@ plainto_tsquery('russian', $1)
						  OR $1 = ''
					ORDER BY o.name ` + params.Desc + `, v.name
					LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(getVacancies, params.Request, 10, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]*baseModels.Vacancy, 0)

	for rows.Next() {
		var vacancyDB baseModels.Vacancy
		err := rows.Scan(&vacancyDB.Organization.ID, &vacancyDB.Organization.Avatar, &vacancyDB.Organization.Name, &vacancyDB.ID,
			&vacancyDB.Name, &vacancyDB.Keywords, &vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax)
		if err != nil {
			return nil, err
		}

		result = append(result, &vacancyDB)
	}

	return result, nil
}
