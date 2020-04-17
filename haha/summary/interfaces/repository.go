package summaryInterfaces

import (
	"joblessness/haha/models/base"
)

type SummaryRepository interface {
	CreateSummary(summary *baseModels.Summary) (summaryID uint64, err error)
	GetAllSummaries(page int) (summaries baseModels.Summaries, err error)
	GetUserSummaries(page int, userID uint64) (summaries baseModels.Summaries, err error)
	GetSummary(summaryID uint64) (summary *baseModels.Summary, err error)
	CheckAuthor(summaryID uint64, authorID uint64) (err error)
	ChangeSummary(summary *baseModels.Summary) (err error)
	DeleteSummary(summaryID uint64) (err error)
	SendSummary(sendSummary *baseModels.SendSummary) (err error)
	RefreshSummary(summaryID, vacancyID uint64) (err error)
	IsOrganizationVacancy(vacancyID, userID uint64) (err error)
	ResponseSummary(sendSummary *baseModels.SendSummary) (err error)
	GetOrgSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error)
	GetUserSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error)
	SendSummaryByMail(summaryID uint64, to string) (err error)
}
