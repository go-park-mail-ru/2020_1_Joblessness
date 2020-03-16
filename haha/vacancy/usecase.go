package vacancy

import "joblessness/haha/models"

type VacancyUseCase interface {
	CreateVacancy(vacancy models.Vacancy) (uint64, error)
	GetVacancies() ([]models.Vacancy, error)
	GetVacancy(vacancyID uint64) (models.Vacancy, error)
	ChangeVacancy(vacancy models.Vacancy) error
	DeleteVacancy(vacancyID uint64) error
}

