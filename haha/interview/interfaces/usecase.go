package interviewInterfaces

import (
	"joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
)

type InterviewUseCase interface {
	ResponseSummary(sendSummary *baseModels.SendSummary) (err error)
	SaveMessage(message *chat.Message) (err error)
	GetHistory(parameters *baseModels.ChatParameters) (result baseModels.Messages, err error)
	GetResponseCredentials(summaryID, vacancyID uint64) (result *baseModels.SummaryCredentials, err error)
	GetConversations(userID uint64) (result baseModels.Conversations, err error)
}
