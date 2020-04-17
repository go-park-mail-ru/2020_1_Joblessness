package summaryInterfaces

import (
	"joblessness/haha/models/base"
)

type SummaryUseCase interface {
	CreateSummary(summary *baseModels.Summary) (summaryID uint64, err error)
	GetAllSummaries(page string) (summaries baseModels.Summaries, err error)
	GetUserSummaries(page string, userID uint64) (summaries baseModels.Summaries, err error)
	GetSummary(summaryID uint64) (summary *baseModels.Summary, err error)
	ChangeSummary(summary *baseModels.Summary) (err error)
	DeleteSummary(summaryID, authorID uint64) (err error)
	SendSummary(sendSummary *baseModels.SendSummary) (err error)
	ResponseSummary(sendSummary *baseModels.SendSummary) (err error)
	GetOrgSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error)
	GetUserSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error)
	SendSummaryByMail(summaryID, authorID uint64, to string) (err error)
}
