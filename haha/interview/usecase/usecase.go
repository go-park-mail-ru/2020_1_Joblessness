package interviewUseCase

import (
	"github.com/microcosm-cc/bluemonday"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	"joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
)

type InterviewUseCase struct {
	interviewRepo interviewInterfaces.InterviewRepository
	policy        *bluemonday.Policy
}

func NewInterviewUseCase(interviewRepo interviewInterfaces.InterviewRepository, policy *bluemonday.Policy) *InterviewUseCase {
	return &InterviewUseCase{
		interviewRepo: interviewRepo,
		policy:        policy,
	}
}

func (u *InterviewUseCase) ResponseSummary(sendSummary *baseModels.SendSummary) (err error) {
	err = u.interviewRepo.IsOrganizationVacancy(sendSummary.VacancyID, sendSummary.OrganizationID)
	if err != nil {
		return err
	}

	err = u.interviewRepo.ResponseSummary(sendSummary)

	return err
}

func (u *InterviewUseCase) SaveMessage(message *chat.Message) (err error) {
	return u.interviewRepo.SaveMessage(message)
}

func (u *InterviewUseCase) GetHistory(parameters *baseModels.ChatParameters) (result baseModels.Messages, err error) {
	return u.interviewRepo.GetHistory(parameters)
}

func (u *InterviewUseCase) GetResponseCredentials(summaryID, vacancyID uint64) (result *baseModels.SummaryCredentials, err error) {
	return u.interviewRepo.GetResponseCredentials(summaryID, vacancyID)
}

func (u *InterviewUseCase) GetConversations(userID uint64) (result baseModels.Conversations, err error) {
	return u.interviewRepo.GetConversations(userID)
}
