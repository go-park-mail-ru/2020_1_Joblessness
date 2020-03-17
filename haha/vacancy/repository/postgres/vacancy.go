package postgres

import (
	"database/sql"
	"joblessness/haha/models"
	"time"
)

type Vacancy struct {
	ID uint64
	OrganizationID uint64
	Name sql.NullString
	Description sql.NullString
	SalaryFrom int
	SalaryTo int
	WithTax bool
	Responsibilities sql.NullString
	Conditions sql.NullString
	Keywords sql.NullString
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

type Organization struct {
	ID uint64
	Name string
	Site string
}

func toPostgres(v *models.Vacancy) *Vacancy {
	return &Vacancy{
		ID:               v.ID,
		OrganizationID:   v.Organization.ID,
		Name:             sql.NullString{String: v.Name, Valid: true},
		Description:      sql.NullString{String: v.Description, Valid: true},
		SalaryFrom:       v.SalaryFrom,
		SalaryTo:         v.SalaryTo,
		WithTax:          v.WithTax,
		Responsibilities: sql.NullString{String: v.Responsibilities, Valid: true},
		Conditions:       sql.NullString{String: v.Conditions, Valid: true},
		Keywords:         sql.NullString{String: v.Keywords, Valid: true},
	}
}

func toModel(v *Vacancy, u *User, o *Organization) *models.Vacancy {
	organization := models.VacancyOrganization{
		ID:     u.ID,
		Tag:    u.Tag,
		Email:  u.Email,
		Phone:  u.Phone,
		Avatar: u.Avatar,
		Name:   o.Name,
		Site:   o.Site,
	}

	return &models.Vacancy{
		ID:               v.ID,
		Organization:     organization,
		Name:             v.Name.String,
		Description:      v.Description.String,
		SalaryFrom:       v.SalaryFrom,
		SalaryTo:         v.SalaryTo,
		WithTax:          v.WithTax,
		Responsibilities: v.Responsibilities.String,
		Conditions:       v.Conditions.String,
		Keywords:         v.Keywords.String,
	}
}

type VacancyRepository struct {
	db *sql.DB
}

func NewVacancyRepository(db *sql.DB) *VacancyRepository {
	return &VacancyRepository{db}
}

func (r *VacancyRepository) CreateVacancy(vacancy *models.Vacancy) (vacancyID uint64, err error) {
	vacancyDB := toPostgres(vacancy)

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

func (r *VacancyRepository) GetVacancyOrganization(organizationID uint64) (*User, *Organization, error) {
	user := User{ID: organizationID}

	getUser := `SELECT person_id, tag, email, phone, avatar
				FROM users WHERE id = $1`
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.PersonID, &user.Tag, &user.Email, &user.Phone, &user.Avatar)
	if err != nil {
		return nil, nil, err
	}

	var organization Organization

	getOrganization := `SELECT name, site
						FROM organization WHERE id = $1`
	err = r.db.QueryRow(getOrganization, user.OrganizationID).Scan(&organization.Name, &organization.Site)
	if err != nil {
		return nil, nil, err
	}

	return &user, &organization, nil
}

func (r *VacancyRepository) GetVacancies() (vacancies []models.Vacancy, err error) {
	getVacancies := `SELECT id, organization_id, name, description, salary_from, salary_to, with_tax, responsibilities,
       						conditions, keywords
					 FROM vacancy;`
	rows, err := r.db.Query(getVacancies)
	if err != nil {
		return vacancies, err
	}

	for rows.Next() {
		var vacancyDB Vacancy
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

		vacancies = append(vacancies, *toModel(&vacancyDB, userDB, organizationDB))
	}

	return vacancies, nil
}

func (r *VacancyRepository) GetVacancy(vacancyID uint64) (vacancy *models.Vacancy, err error) {
	var vacancyDB Vacancy

	getVacancy := `SELECT id, organization_id, name, description, salary_from, salary_to, with_tax, responsibilities,
       					  conditions, keywords
				   FROM vacancy WHERE id = $1;`
	err = r.db.QueryRow(getVacancy, vacancyID).Scan(vacancyDB)
	if err != nil {
		return vacancy, err
	}

	userDB, organizationDB, err := r.GetVacancyOrganization(vacancyDB.OrganizationID)
	if err != nil {
		return nil, err
	}

	return toModel(&vacancyDB, userDB, organizationDB), nil
}

func (r *VacancyRepository) ChangeVacancy(vacancy *models.Vacancy) (err error) {
	vacancyDB := toPostgres(vacancy)
	//TODO проверять автора

	changeVacancy := `UPDATE vacancy
					  SET organization_id = $1, name = $2, description = $3, salary_from = $4, salary_to = $5,
						  with_tax = $6, responsibilities = $7, conditions = $8, keywords = $9
					  WHERE id = $10;`
	_, err = r.db.Exec(changeVacancy, &vacancyDB.OrganizationID, &vacancyDB.Name, &vacancyDB.Description,
						 &vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax, &vacancyDB.Responsibilities,
						 &vacancyDB.Conditions, &vacancyDB.Keywords, &vacancyDB.ID)
	if err != nil {
		return err
	}

	return nil
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
