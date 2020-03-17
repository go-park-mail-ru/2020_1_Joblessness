package usecase

import (
	"github.com/kataras/golog"
	"joblessness/haha/models"
	"joblessness/haha/summary"
)

type SummaryUseCase struct {
	summaryRepo summary.SummaryRepository
}

func NewSummaryUseCase(summaryRepo summary.SummaryRepository) *SummaryUseCase {
	return &SummaryUseCase{summaryRepo}
}

func (u *SummaryUseCase) CreateSummary(summary *models.Summary) (summaryID uint64, err error) {
	golog.Info("hello")
	return u.summaryRepo.CreateSummary(summary)
}

func (u *SummaryUseCase) GetAllSummaries() (summaries []models.Summary, err error) {
	return u.summaryRepo.GetAllSummaries()
}

func (u *SummaryUseCase) GetUserSummaries(userID uint64) (summaries []models.Summary, err error) {
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
