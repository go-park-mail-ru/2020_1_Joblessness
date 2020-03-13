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

	createVacancy := `INSERT INTO vacancy (organization_id, name, description, salary_from, salary_to, with_tax, responsibilities, conditions, keywords)
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	err = r.db.QueryRow(createVacancy, vacancyDB.OrganizationID, vacancyDB.Name, vacancyDB.Description, vacancyDB.SalaryFrom, vacancyDB.SalaryTo, vacancyDB.WithTax, vacancyDB.Responsibilities, vacancyDB.Conditions, vacancyDB.Keywords).Scan(&vacancyID)
	if err != nil {
		return 0, err
	}

	createRequirements := `INSERT INTO requirements (vacancy_id, driver_license, has_car, schedule, employment)
						 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.db.Exec(createRequirements, requirementsDB.VacancyID, requirementsDB.DriverLicense, requirementsDB.HasCar, requirementsDB.Schedule, requirementsDB.Employment)
	if err != nil {
		return 0, err
	}

	return vacancyID, nil
}

func (r *VacancyRepository) GetVacancies() ([]models.Vacancy, error) {
	// TODO: implementation
	return []models.Vacancy{}, nil
}

func (r *VacancyRepository) GetVacancy(vacancyID int) (models.Vacancy, error) {
	// TODO: implementation
	return models.Vacancy{}, nil
}

func (r *VacancyRepository) ChangeVacancy(vacancy models.Vacancy) error {
	// TODO: implementation
	return nil
}

func (r *VacancyRepository) DeleteVacancy(vacancyID int) error {
	// TODO: implementation
	return nil
}
