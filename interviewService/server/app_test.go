package interviewServer

//go:generate mockgen -destination=../../haha/interview/repository/mock/interview.go -package=mock joblessness/haha/interview/interfaces InterviewRepository

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	interviewGrpc "joblessness/haha/interview/repository/grpc"
	"joblessness/haha/interview/repository/mock"
	baseModels "joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
	interviewRpc "joblessness/interviewService/rpc"
	"net"
	"testing"
)

type userSuite struct {
	suite.Suite
	controller *gomock.Controller
	grpcRepo   *interviewGrpc.InterviewGrpcRepository
	repo       *mock.MockInterviewRepository
	server     *grpc.Server
	list       net.Listener
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
	interviewConn, err := grpc.Dial(
		"127.0.0.1:8003",
		grpc.WithInsecure(),
	)
	assert.NoError(suite.T(), err, "Unable to start server")

	suite.grpcRepo = interviewGrpc.NewInterviewGrpcRepository(interviewConn)
	assert.NoError(suite.T(), err)

	suite.repo = mock.NewMockInterviewRepository(suite.controller)
	suite.list, err = net.Listen("tcp", "127.0.0.1:8003")
	assert.NoError(suite.T(), err, "Unable to listen")
	suite.server = grpc.NewServer()
	interviewRpc.RegisterInterviewServer(suite.server, NewInterviewServer(suite.repo))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestIsOrganizationVacancy() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().IsOrganizationVacancy(uint64(1), uint64(2)).Times(1).Return(nil)

	err := suite.grpcRepo.IsOrganizationVacancy(uint64(1), uint64(2))

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestIsOrganizationVacancyNotOwner() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().IsOrganizationVacancy(uint64(1), uint64(2)).Times(1).Return(interviewInterfaces.ErrOrganizationIsNotOwner)

	err := suite.grpcRepo.IsOrganizationVacancy(uint64(1), uint64(2))

	assert.Equal(suite.T(), interviewInterfaces.ErrOrganizationIsNotOwner, err)
}

func (suite *userSuite) TestIsOrganizationVacancyDefaultErr() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().IsOrganizationVacancy(uint64(1), uint64(2)).Times(1).Return(errors.New(""))

	err := suite.grpcRepo.IsOrganizationVacancy(uint64(1), uint64(2))

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestResponseSummary() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().ResponseSummary(gomock.Any()).Times(1).Return(nil)

	err := suite.grpcRepo.ResponseSummary(&baseModels.SendSummary{})

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSaveMessage() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().SaveMessage(gomock.Any()).Times(1).Return(nil)

	err := suite.grpcRepo.SaveMessage(&chat.Message{})

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetHistory() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().GetHistory(gomock.Any()).Times(1).Return(baseModels.Messages{
		From: []*chat.Message{&chat.Message{}}, To: []*chat.Message{&chat.Message{}}}, nil)

	_, err := suite.grpcRepo.GetHistory(&baseModels.ChatParameters{})

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetResponseCredentials() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().GetResponseCredentials(uint64(1), uint64(2)).Times(1).Return(&baseModels.SummaryCredentials{}, nil)

	_, err := suite.grpcRepo.GetResponseCredentials(uint64(1), uint64(2))

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetConversations() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().GetConversations(uint64(1)).Times(1).Return(baseModels.Conversations{&baseModels.ConversationTitle{}}, nil)

	_, err := suite.grpcRepo.GetConversations(uint64(1))

	assert.NoError(suite.T(), err)
}
