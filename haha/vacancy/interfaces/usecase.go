package vacancyInterfaces

import (
	"joblessness/haha/models/base"
)

type VacancyUseCase interface {
	CreateVacancy(vacancy *baseModels.Vacancy) (uint64, error)
	GetVacancies(page string) (baseModels.Vacancies, error)
	GetVacancy(vacancyID uint64) (*baseModels.Vacancy, error)
	ChangeVacancy(vacancy *baseModels.Vacancy) error
	DeleteVacancy(vacancyID uint64, authorID uint64) error
	GetOrgVacancies(userID uint64) (baseModels.Vacancies, error)
}
