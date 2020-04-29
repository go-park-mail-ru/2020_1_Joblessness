package interviewUseCase

import (
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/interview/repository/mock"
	baseModels "joblessness/haha/models/base"
	mockRoom "joblessness/haha/utils/chat/mock"
	"testing"
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
