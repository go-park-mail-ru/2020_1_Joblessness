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

type Requirements struct {
	ID uint64
	VacancyID uint64
	DriverLicense string
	HasCar bool
	Schedule string
	Employment string
}

func toPostgres(v *models.Vacancy) (*Vacancy, *Requirements) {
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
	},
	&Requirements{
		ID:            0,
		VacancyID:     0,
		DriverLicense: "",
		HasCar:        false,
		Schedule:      "",
		Employment:    "",
	}
}

func toModel(v *Vacancy, r *Requirements) *models.Vacancy {
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
	vacancyDB, requirementsDB := toPostgres(&vacancy)

	createVacancy := `INSERT INTO vacancy (organization_id, name, description, salary_from, salary_to, with_tax,
                     					   responsibilities, conditions, keywords)
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	err = r.db.QueryRow(createVacancy, vacancyDB.OrganizationID, vacancyDB.Name, vacancyDB.Description,
						vacancyDB.SalaryFrom, vacancyDB.SalaryTo, vacancyDB.WithTax, vacancyDB.Responsibilities,
						vacancyDB.Conditions, vacancyDB.Keywords).Scan(&vacancyID)
	if err != nil {
		return vacancyID, err
	}

	createRequirements := `INSERT INTO requirements (vacancy_id, driver_license, has_car, schedule, employment)
						 VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = r.db.Exec(createRequirements, requirementsDB.VacancyID, requirementsDB.DriverLicense,
					   requirementsDB.HasCar, requirementsDB.Schedule, requirementsDB.Employment)
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

	getRequirements := `SELECT id, driver_license, has_car, schedule, employment
						FROM requirements WHERE vacancy_id = $1;`

	for rows.Next() {
		var vacancyDB Vacancy
		err = rows.Scan(&vacancyDB.OrganizationID, &vacancyDB.Name, &vacancyDB.Description, &vacancyDB.SalaryFrom,
						&vacancyDB.SalaryTo, &vacancyDB.WithTax, &vacancyDB.Responsibilities, &vacancyDB.Conditions,
						&vacancyDB.Keywords)
		if err != nil {
			return vacancies, err
		}

		var requirementsDB Requirements

		err = r.db.QueryRow(getRequirements, vacancyDB.ID).Scan(&requirementsDB)
		if err != nil {
			return vacancies, err
		}

		vacancies = append(vacancies, *toModel(&vacancyDB, &requirementsDB))
	}

	return vacancies, nil
}

func (r *VacancyRepository) GetVacancy(vacancyID int) (vacancy models.Vacancy, err error) {
	var vacancyDB Vacancy

	getVacancy := `SELECT id, organization_id, name, description, salary_from, salary_to, with_tax, responsibilities,
       					  conditions, keywords
				   FROM vacancy WHERE id = $1;`
	err = r.db.QueryRow(getVacancy, vacancyID).Scan(vacancyDB)
	if err != nil {
		return vacancy, err
	}

	var requirementsDB Requirements

	getRequirements := `SELECT id, driver_license, has_car, schedule, employment
						FROM requirements WHERE vacancy_id = $1;`
	err = r.db.QueryRow(getRequirements, vacancyDB.ID).Scan(&requirementsDB)
	if err != nil {
		return vacancy, err
	}

	return *toModel(&vacancyDB, &requirementsDB), nil
}

func (r *VacancyRepository) ChangeVacancy(vacancy models.Vacancy) (err error) {
	vacancyDB, requirementsDB := toPostgres(&vacancy)

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

	changeRequirements := `UPDATE requirements
						   SET driver_license = $1, has_car = $2, schedule = $3, employment = $4
						   WHERE vacancy_id = $5`
	_, err = r.db.Exec(changeRequirements, &requirementsDB.DriverLicense, &requirementsDB.HasCar,
					   &requirementsDB.Schedule, &requirementsDB.Employment, &requirementsDB.VacancyID)
	if err != nil {
		return err
	}

	return nil
}

func (r *VacancyRepository) DeleteVacancy(vacancyID int) (err error) {
	deleteVacancy := `DELETE FROM vacancy
					  WHERE id = $1`
	_, err = r.db.Exec(deleteVacancy, vacancyID)
	if err != nil {
		return err
	}

	return nil
}
