package vacancy

import "joblessness/haha/models"

type UseCase interface {
	CreateVacancy(vacancy models.Vacancy) (uint64, error)
	GetVacancies() ([]models.Vacancy, error)
	GetVacancy(vacancyID uint64) (models.Vacancy, error)
	ChangeVacancy(vacancy models.Vacancy) error
	DeleteVacancy(vacancyID uint64) error
}

func (u UseCase) ChangeSummary(summary *models.Summary) (err error) {
	panic("implement me")
}

func (u UseCase) CreateSummary(summary *models.Summary) (summaryID uint64, err error) {
	panic("implement me")
}

func (u UseCase) DeleteSummary(summaryID uint64) (err error) {
	panic("implement me")
}

func (u UseCase) GetSummaries() (summaries *[]models.Summary, err error) {
	panic("implement me")
}

func (u UseCase) GetSummary(summaryID uint64) (summary *models.Summary, err error) {
	panic("implement me")
}

