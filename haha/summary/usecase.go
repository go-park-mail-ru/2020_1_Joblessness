package summary

import "joblessness/haha/models"

type SummaryUseCase interface {
	CreateSummary(summary *models.Summary) (summaryID uint64, err error)
	GetAllSummaries(pageNumber uint64) (summaries []models.Summary, pageCount uint64, hasPrev, hasNext bool, err error)
	GetUserSummaries(userID uint64, pageNumber uint64) (summaries []models.Summary, pageCount uint64, hasPrev, hasNext bool, err error)
	GetSummary(summaryID uint64) (summary *models.Summary, err error)
	ChangeSummary(summary *models.Summary) (err error)
	DeleteSummary(summaryID uint64) (err error)
}
