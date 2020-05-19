package interviewUseCase

import (
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/interview/repository/mock"
	baseModels "joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
	mockRoom "joblessness/haha/utils/chat/mock"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	controller   *gomock.Controller
	uc           *InterviewUseCase
	person       baseModels.Person
	organization baseModels.Organization
	repo         *mock.MockInterviewRepository
	sidEx        string
	room         *mockRoom.MockRoom
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
	defer suite.controller.Finish()

	suite.repo = mock.NewMockInterviewRepository(suite.controller)
	suite.uc = NewInterviewUseCase(suite.repo, bluemonday.UGCPolicy())
	suite.room = mockRoom.NewMockRoom(suite.controller)

	//suite.room.EXPECT().Run().Times(1)
	//suite.uc.EnableRoom(suite.room)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestGenerateMessage() {
	sendSummary := &baseModels.SendSummary{}
	credits := &baseModels.SummaryCredentials{}
	suite.repo.EXPECT().GetResponseCredentials(sendSummary.SummaryID, sendSummary.VacancyID).Return(credits, nil).Times(1)

	_, err := suite.uc.generateMessage(sendSummary)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSaveMessage() {
	message := chat.Message{
		Message:   "message",
		UserOneID: 1,
		UserOne:   "awd",
		UserTwoID: 2,
		UserTwo:   "",
		Created:   time.Now(),
		VacancyID: 0,
	}
	suite.repo.EXPECT().SaveMessage(&message).Return(nil).Times(1)

	err := suite.uc.SaveMessage(&message)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetHistory() {
	params := baseModels.ChatParameters{
		From: 1,
		To:   2,
		Page: 3,
	}
	suite.repo.EXPECT().GetHistory(&params).Return(baseModels.Messages{}, nil).Times(1)

	_, err := suite.uc.GetHistory(&params)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetConversations() {
	suite.repo.EXPECT().GetConversations(uint64(2)).Return(baseModels.Conversations{}, nil).Times(1)

	_, err := suite.uc.GetConversations(uint64(2))

	assert.NoError(suite.T(), err)
}
