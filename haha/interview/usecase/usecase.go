package interviewUseCase

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	"joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
	"time"
)

const (
	messageBufferSize = 256
)

type InterviewUseCase struct {
	interviewRepo interviewInterfaces.InterviewRepository
	room          chat.Room
	policy        *bluemonday.Policy
}

func NewInterviewUseCase(interviewRepo interviewInterfaces.InterviewRepository,
	policy *bluemonday.Policy) (useCase *InterviewUseCase) {
	useCase = &InterviewUseCase{
		interviewRepo: interviewRepo,
		policy:        policy,
	}
	return useCase
}

func (u *InterviewUseCase) EnableRoom(room chat.Room) {
	u.room = room
	go u.room.Run()
}

func (u *InterviewUseCase) EnterChat(userID uint64, socket *websocket.Conn) {
	chatter := &chat.Chatter{
		ID:     userID,
		Socket: socket,
		Send:   make(chan []byte, messageBufferSize),
		Room:   u.room,
	}

	u.room.Join(chatter)
	defer func() {
		u.room.Leave(chatter)
	}()
	go chatter.Write()
	chatter.Read()
}

func (u *InterviewUseCase) generateMessage(sendSummary *baseModels.SendSummary) (result *chat.Message, err error) {
	credentials, err := u.GetResponseCredentials(sendSummary.SummaryID, sendSummary.VacancyID)
	if err != nil {
		return nil, err
	}

	var status string
	if sendSummary.Accepted {
		status = "одобрено."
	} else if sendSummary.Denied {
		status = "отклонено."
	} else {
		status = "просмотренно, Вы приглашены на собеседование."
	}

	return &chat.Message{
		Message:   fmt.Sprintf("Ваше резюме было %s", status),
		UserOneID: credentials.OrganizationID,
		UserOne:   credentials.OrganizationName,
		UserTwoID: credentials.UserID,
		Created:   time.Now(),
	}, nil
}

func (u *InterviewUseCase) ResponseSummary(sendSummary *baseModels.SendSummary) (err error) {
	err = u.interviewRepo.IsOrganizationVacancy(sendSummary.VacancyID, sendSummary.OrganizationID)
	if err != nil {
		return err
	}
	err = u.interviewRepo.ResponseSummary(sendSummary)
	if err == nil {
		message, err := u.generateMessage(sendSummary)
		if err != nil {
			return err
		}
		if message != nil {
			err = u.room.SendGeneratedMessage(message)
			if err != nil {
				return err
			}
		}
	}

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
