package summaryInterfaces

import "joblessness/haha/models"

type SummaryRepository interface {
	CreateSummary(summary *models.Summary) (summaryID uint64, err error)
	GetAllSummaries(page int) (summaries []models.Summary, err error)
	GetUserSummaries(userID uint64) (summaries []models.Summary, err error)
	GetSummary(summaryID uint64) (summary *models.Summary, err error)
	ChangeSummary(summary *models.Summary) (err error)
	DeleteSummary(summaryID uint64) (err error)
	IsPersonSummary(summaryID, userID uint64) (res bool, err error)
	SendSummary(sendSummary *models.SendSummary) (err error)
	RefreshSummary(summaryID, vacancyID uint64) (err error)
	IsOrganizationVacancy(vacancyID, userID uint64) (res bool, err error)
	ResponseSummary(sendSummary *models.SendSummary) (err error)
	GetOrgSummaries(userID uint64) (summaries models.OrgSummaries, err error)
}
