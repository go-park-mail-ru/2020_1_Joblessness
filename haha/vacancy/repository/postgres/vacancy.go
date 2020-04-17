package vacancyPostgres

import (
	"database/sql"
	"joblessness/haha/models/base"
	"joblessness/haha/models/postgres"
	"joblessness/haha/vacancy/interfaces"
)

type VacancyRepository struct {
	db *sql.DB
}

func NewVacancyRepository(db *sql.DB) *VacancyRepository {
	return &VacancyRepository{db}
}

func (r *VacancyRepository) CreateVacancy(vacancy *baseModels.Vacancy) (vacancyID uint64, err error) {
	vacancyDB := pgModels.ToPgVacancy(vacancy)

	createVacancy := `INSERT INTO vacancy (organization_id, name, description, salary_from, salary_to, with_tax,
                     					   responsibilities, conditions, keywords)
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	err = r.db.QueryRow(createVacancy, &vacancyDB.OrganizationID, &vacancyDB.Name, &vacancyDB.Description,
		&vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax, &vacancyDB.Responsibilities,
		&vacancyDB.Conditions, &vacancyDB.Keywords).Scan(&vacancyID)
	if err != nil {
		return vacancyID, err
	}

	return vacancyID, nil
}

func (r *VacancyRepository) GetVacancyOrganization(organizationID uint64) (*pgModels.User, *pgModels.Organization, error) {
	user := pgModels.User{ID: organizationID}
	var organization pgModels.Organization

	getUser := `SELECT organization_id, tag, email, phone, avatar, name, site
				FROM users u
				JOIN organization o on u.organization_id = o.id
				WHERE u.id = $1 AND organization_id IS NOT NULL`
	err := r.db.QueryRow(getUser, user.ID).
		Scan(&user.OrganizationID, &user.Tag, &user.Email, &user.Phone, &user.Avatar, &organization.Name, &organization.Site)
	if err != nil {
		return nil, nil, err
	}

	organization.ID = user.OrganizationID
	return &user, &organization, nil
}

func (r *VacancyRepository) GetVacancy(vacancyID uint64) (vacancy *baseModels.Vacancy, err error) {
	var vacancyDB pgModels.Vacancy

	getVacancy := `SELECT v.id, v.organization_id, v.name, v.description, v.salary_from, v.salary_to, v.with_tax, v.responsibilities,
       					  v.conditions, v.keywords
				   FROM vacancy v WHERE v.id = $1;`
	err = r.db.QueryRow(getVacancy, vacancyID).Scan(&vacancyDB.ID, &vacancyDB.OrganizationID, &vacancyDB.Name,
		&vacancyDB.Description, &vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax, &vacancyDB.Responsibilities,
		&vacancyDB.Conditions, &vacancyDB.Keywords)
	if err != nil {
		return nil, err
	}

	userDB, organizationDB, err := r.GetVacancyOrganization(vacancyDB.OrganizationID)
	if err != nil {
		return nil, err
	}

	return pgModels.ToBaseVacancy(&vacancyDB, userDB, organizationDB), nil
}

func (r *VacancyRepository) GetVacancies(page int) (vacancies baseModels.Vacancies, err error) {
	getVacancies := `SELECT id, organization_id, name, description, salary_from, salary_to, with_tax, responsibilities,
       						conditions, keywords
					 FROM vacancy
					LIMIT $1 OFFSET $2;`
	rows, err := r.db.Query(getVacancies, 10, page*10)
	if err != nil {
		return vacancies, err
	}
	defer rows.Close()

	vacancies = make(baseModels.Vacancies, 0)

	for rows.Next() {
		var vacancyDB pgModels.Vacancy
		err = rows.Scan(&vacancyDB.ID, &vacancyDB.OrganizationID, &vacancyDB.Name, &vacancyDB.Description,
			&vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax, &vacancyDB.Responsibilities,
			&vacancyDB.Conditions, &vacancyDB.Keywords)
		if err != nil {
			return vacancies, err
		}

		userDB, organizationDB, err := r.GetVacancyOrganization(vacancyDB.OrganizationID)
		if err != nil {
			return nil, err
		}

		vacancies = append(vacancies, pgModels.ToBaseVacancy(&vacancyDB, userDB, organizationDB))
	}

	return vacancies, nil
}

func (r *VacancyRepository) CheckAuthor(vacancyID, authorID uint64) (err error) {
	var isAuthor bool

	checkAuthor := `SELECT organization_id = $1 FROM vacancy WHERE id = $2`
	if err = r.db.QueryRow(checkAuthor, authorID, vacancyID).Scan(&isAuthor); err != nil {
		return vacancyInterfaces.ErrOrgIsNotOwner
	}

	return err
}

func (r *VacancyRepository) ChangeVacancy(vacancy *baseModels.Vacancy) (err error) {
	vacancyDB := pgModels.ToPgVacancy(vacancy)

	changeVacancy := `UPDATE vacancy
					  SET name = $1, description = $2, salary_from = $3, salary_to = $4,
						  with_tax = $5, responsibilities = $6, conditions = $7, keywords = $8
					  WHERE id = $9;`
	_, err = r.db.Exec(changeVacancy, vacancyDB.Name, vacancyDB.Description,
		vacancyDB.SalaryFrom, vacancyDB.SalaryTo, vacancyDB.WithTax, vacancyDB.Responsibilities,
		vacancyDB.Conditions, vacancyDB.Keywords, vacancyDB.ID)

	return err
}

func (r *VacancyRepository) DeleteVacancy(vacancyID uint64) (err error) {
	//TODO проверять автора
	deleteVacancy := `DELETE FROM vacancy
					  WHERE id = $1`
	_, err = r.db.Exec(deleteVacancy, vacancyID)
	if err != nil {
		return err
	}

	return nil
}

func (r *VacancyRepository) GetOrgVacancies(userID uint64) (vacancies baseModels.Vacancies, err error) {
	getVacancies := `SELECT id, name, salary_from, salary_to, with_tax, keywords
					 FROM vacancy
					WHERE organization_id = $1;`
	rows, err := r.db.Query(getVacancies, userID)
	if err != nil {
		return vacancies, err
	}
	defer rows.Close()

	vacancies = make(baseModels.Vacancies, 0)

	for rows.Next() {
		var vacancyDB baseModels.Vacancy
		err = rows.Scan(&vacancyDB.ID, &vacancyDB.Name, &vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax,
			&vacancyDB.Keywords)
		if err != nil {
			return vacancies, err
		}

		vacancies = append(vacancies, &vacancyDB)
	}

	return vacancies, nil
}
