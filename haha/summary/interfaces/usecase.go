package summaryInterfaces

import "joblessness/haha/models"

type SummaryUseCase interface {
	CreateSummary(summary *models.Summary) (summaryID uint64, err error)
	GetAllSummaries() (summaries []models.Summary, err error)
	GetUserSummaries(userID uint64) (summaries []models.Summary, err error)
	GetSummary(summaryID uint64) (summary *models.Summary, err error)
	ChangeSummary(summary *models.Summary) (err error)
	DeleteSummary(summaryID uint64) (err error)
	SendSummary(sendSummary *models.SendSummary) (err error)
	ResponseSummary(sendSummary *models.SendSummary) (err error)
	GetOrgSummaries(userID uint64) (summaries models.OrgSummaries, err error)
}
