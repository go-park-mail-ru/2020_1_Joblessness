package postgres

import (
	"database/sql"
	"joblessness/haha/models"
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

func toPostgres(v *models.Vacancy) *Vacancy {
	return &Vacancy{
		ID:               v.ID,
		OrganizationID:   0,
		Name:             v.Name,
		Description:      v.Description,
		SalaryFrom:       v.Salary,
		SalaryTo:         0,
		WithTax:          false,
		Responsibilities: "",
		Conditions:       "",
		Keywords:         "",
	}
}

func toModel(v *Vacancy) *models.Vacancy {
	return &models.Vacancy{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
		Skills:      "",
		Salary:      v.SalaryFrom,
		Address:     "",
		PhoneNumber: "",
	}
}

type VacancyRepository struct {
	db *sql.DB
}

func NewVacancyRepository(db *sql.DB) *VacancyRepository {
	return &VacancyRepository{db}
}

func (r *VacancyRepository) CreateVacancy(vacancy models.Vacancy) (vacancyID uint64, err error) {
	vacancyDB := toPostgres(&vacancy)

	vacancyDB.OrganizationID = 1

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

		vacancies = append(vacancies, *toModel(&vacancyDB))
	}

	return vacancies, nil
}

func (r *VacancyRepository) GetVacancy(vacancyID uint64) (vacancy models.Vacancy, err error) {
	var vacancyDB Vacancy

	getVacancy := `SELECT id, organization_id, name, description, salary_from, salary_to, with_tax, responsibilities,
       					  conditions, keywords
				   FROM vacancy WHERE id = $1;`
	err = r.db.QueryRow(getVacancy, vacancyID).Scan(vacancyDB)
	if err != nil {
		return vacancy, err
	}

	return *toModel(&vacancyDB), nil
}

func (r *VacancyRepository) ChangeVacancy(vacancy models.Vacancy) (err error) {
	vacancyDB := toPostgres(&vacancy)

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
	deleteVacancy := `DELETE FROM vacancy
					  WHERE id = $1`
	_, err = r.db.Exec(deleteVacancy, vacancyID)
	if err != nil {
		return err
	}

	return nil
}
