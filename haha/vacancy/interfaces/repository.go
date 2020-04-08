package vacancyInterfaces

import "joblessness/haha/models"

type VacancyRepository interface {
	CreateVacancy(vacancy *models.Vacancy) (uint64, error)
	GetVacancies(page int) ([]models.Vacancy, error)
	GetVacancy(vacancyID uint64) (*models.Vacancy, error)
	CheckAuthor(vacancyID, authorID uint64) (err error)
	ChangeVacancy(vacancy *models.Vacancy) error
	DeleteVacancy(vacancyID uint64) error
	GetOrgVacancies(userID uint64) ([]models.Vacancy, error)
}
