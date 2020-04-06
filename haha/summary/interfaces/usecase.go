package summaryInterfaces

import "joblessness/haha/models"

type SummaryUseCase interface {
	CreateSummary(summary *models.Summary) (summaryID uint64, err error)
	GetAllSummaries(page string) (summaries []models.Summary, err error)
	GetUserSummaries(page string, userID uint64) (summaries []models.Summary, err error)
	GetSummary(summaryID uint64) (summary *models.Summary, err error)
	ChangeSummary(summary *models.Summary) (err error)
	DeleteSummary(summaryID uint64, authorID uint64) (err error)
	SendSummary(sendSummary *models.SendSummary) (err error)
	ResponseSummary(sendSummary *models.SendSummary) (err error)
	GetOrgSendSummaries(userID uint64) (summaries models.OrgSummaries, err error)
	GetUserSendSummaries(userID uint64) (summaries models.OrgSummaries, err error)
}
