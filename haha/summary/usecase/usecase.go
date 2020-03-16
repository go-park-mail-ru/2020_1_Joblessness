package usecase

import (
	"joblessness/haha/models"
	"joblessness/haha/summary"
)

type SummaryUseCase struct {
	summaryRepo summary.Repository
}

func NewSummaryUseCase(summaryRepo summary.Repository) *SummaryUseCase {
	return &SummaryUseCase{summaryRepo}
}

func (u *SummaryUseCase) CreateSummary(summary *models.Summary) (summaryID uint64, err error) {
	return u.summaryRepo.CreateSummary(summary)
}

func (u *SummaryUseCase) GetAllSummaries() (summaries *[]models.Summary, err error) {
	return u.summaryRepo.GetAllSummaries()
}

func (u *SummaryUseCase) GetUserSummaries(userID uint64) (summaries *[]models.Summary, err error) {
	return u.summaryRepo.GetUserSummaries(userID)
}

func (u *SummaryUseCase) GetSummary(summaryID uint64) (summary *models.Summary, err error) {
	return u.summaryRepo.GetSummary(summaryID)
}

func (u *SummaryUseCase) ChangeSummary(summary *models.Summary) (err error) {
	return u.summaryRepo.ChangeSummary(summary)
}

func (u *SummaryUseCase) DeleteSummary(summaryID uint64) (err error) {
	return u.summaryRepo.DeleteSummary(summaryID)
}
