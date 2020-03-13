package vacancy

import "joblessness/haha/models"

type VacancyRepository interface {
	CreateVacancy(vacancy models.Vacancy) (uint64, error)
	GetVacancies() ([]models.Vacancy, error)
	GetVacancy(vacancyID int) (models.Vacancy, error)
	ChangeVacancy(vacancy models.Vacancy) error
	DeleteVacancy(vacancyID int) error
}
