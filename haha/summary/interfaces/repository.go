package summaryInterfaces

import "joblessness/haha/models"

type SummaryRepository interface {
	CreateSummary(summary *models.Summary) (summaryID uint64, err error)
	GetAllSummaries(page int) (summaries models.Summaries, err error)
	GetUserSummaries(page int, userID uint64) (summaries models.Summaries, err error)
	GetSummary(summaryID uint64) (summary *models.Summary, err error)
	CheckAuthor(summaryID uint64, authorID uint64) (err error)
	ChangeSummary(summary *models.Summary) (err error)
	DeleteSummary(summaryID uint64) (err error)
	SendSummary(sendSummary *models.SendSummary) (err error)
	RefreshSummary(summaryID, vacancyID uint64) (err error)
	GetOrgSendSummaries(userID uint64) (summaries models.OrgSummaries, err error)
	GetUserSendSummaries(userID uint64) (summaries models.OrgSummaries, err error)
	SendSummaryByMail(summaryID uint64, to string) (err error)
}
