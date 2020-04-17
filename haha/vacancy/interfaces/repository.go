package vacancyInterfaces

import (
	"joblessness/haha/models/base"
)

type VacancyRepository interface {
	CreateVacancy(vacancy *baseModels.Vacancy) (uint64, error)
	GetVacancies(page int) (baseModels.Vacancies, error)
	GetVacancy(vacancyID uint64) (*baseModels.Vacancy, error)
	CheckAuthor(vacancyID, authorID uint64) (err error)
	ChangeVacancy(vacancy *baseModels.Vacancy) error
	DeleteVacancy(vacancyID uint64) error
	GetOrgVacancies(userID uint64) (baseModels.Vacancies, error)
}
