package vacancyInterfaces

import (
	"joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
)

type VacancyRepository interface {
	GetRelatedUsers(organizationID uint64) (res []uint64, orgName string, err error)
	CreateVacancy(vacancy *baseModels.Vacancy) (uint64, error)
	GetVacancies(page int) (baseModels.Vacancies, error)
	GetVacancyOrganization(organizationID uint64) (*pgModels.User, *pgModels.Organization, error)
	GetVacancy(vacancyID uint64) (*baseModels.Vacancy, error)
	CheckAuthor(vacancyID, authorID uint64) (err error)
	ChangeVacancy(vacancy *baseModels.Vacancy) error
	DeleteVacancy(vacancyID uint64) error
	GetOrgVacancies(userID uint64) (baseModels.Vacancies, error)
}
