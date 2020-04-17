package interviewInterfaces

import (
	"joblessness/haha/models"
	"joblessness/haha/utils/chat"
)

type InterviewRepository interface {
	IsOrganizationVacancy(vacancyID, userID uint64) (err error)
	ResponseSummary(sendSummary *models.SendSummary) (err error)
	SaveMessage(message *chat.Message) (err error)
	GetHistory(parameters *models.ChatParameters) (result models.Messages, err error)
	GetResponseCredentials(summaryID, vacancyID uint64) (result *models.SummaryCredentials, err error)
	GetConversations(userID uint64) (result models.Conversations, err error)
}
