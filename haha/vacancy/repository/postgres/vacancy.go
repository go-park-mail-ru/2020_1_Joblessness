package vacancyRepoPostgres

import (
	"database/sql"
	"joblessness/haha/models"
	"time"
)

type Vacancy struct {
	ID uint64
	OrganizationID uint64
	Name string
	Description string
	SalaryFrom int
	SalaryTo int
	WithTax bool
	Responsibilities string
	Conditions string
	Keywords string
}

type User struct {
	ID             uint64
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
		Name:             v.Name,
		Description:      v.Description,
		SalaryFrom:       v.SalaryFrom,
		SalaryTo:         v.SalaryTo,
		WithTax:          v.WithTax,
		Responsibilities: v.Responsibilities,
		Conditions:       v.Conditions,
		Keywords:         v.Keywords,
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
		Name:             v.Name,
		Description:      v.Description,
		SalaryFrom:       v.SalaryFrom,
		SalaryTo:         v.SalaryTo,
		WithTax:          v.WithTax,
		Responsibilities: v.Responsibilities,
		Conditions:       v.Conditions,
		Keywords:         v.Keywords,
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
	var organization Organization

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

func (r *VacancyRepository) GetVacancy(vacancyID uint64) (vacancy *models.Vacancy, err error) {
	var vacancyDB Vacancy

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

	return toModel(&vacancyDB, userDB, organizationDB), nil
}

func (r *VacancyRepository) GetVacancies(page int) (vacancies []models.Vacancy, err error) {
	getVacancies := `SELECT id, organization_id, name, description, salary_from, salary_to, with_tax, responsibilities,
       						conditions, keywords
					 FROM vacancy
					LIMIT $1 OFFSET $2;`
	rows, err := r.db.Query(getVacancies, 10, page*10)
	if err != nil {
		return vacancies, err
	}
	defer rows.Close()

	vacancies = make([]models.Vacancy, 0)

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

func (r *VacancyRepository) ChangeVacancy(vacancy *models.Vacancy) (err error) {
	vacancyDB := toPostgres(vacancy)

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

func (r *VacancyRepository) GetOrgVacancies(userID uint64) (vacancies []models.Vacancy, err error) {
	getVacancies := `SELECT id, name, salary_from, salary_to, with_tax, keywords
					 FROM vacancy
					WHERE organization_id = $1;`
	rows, err := r.db.Query(getVacancies, userID)
	if err != nil {
		return vacancies, err
	}
	defer rows.Close()

	vacancies = make([]models.Vacancy, 0)

	for rows.Next() {
		var vacancyDB models.Vacancy
		err = rows.Scan(&vacancyDB.ID, &vacancyDB.Name, &vacancyDB.SalaryFrom, &vacancyDB.SalaryTo, &vacancyDB.WithTax,
			&vacancyDB.Keywords)
		if err != nil {
			return vacancies, err
		}

		vacancies = append(vacancies, vacancyDB)
	}

	return vacancies, nil
}