package searchServer

//go:generate mockgen -destination=../../haha/search/repository/mock/search.go -package=mock joblessness/haha/search/interfaces SearchRepository

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	baseModels "joblessness/haha/models/base"
	searchGrpc "joblessness/haha/search/repository/grpc"
	"joblessness/haha/search/repository/mock"
	searchRpc "joblessness/searchService/rpc"
	"net"
	"testing"
)

type userSuite struct {
	suite.Suite
	controller *gomock.Controller
	grpcRepo   *searchGrpc.SearchGrpcRepository
	repo 	*mock.MockSearchRepository
	server     *grpc.Server
	list net.Listener
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
	interviewConn, err := grpc.Dial(
		"127.0.0.1:8003",
		grpc.WithInsecure(),
	)
	assert.NoError(suite.T(), err)

	suite.grpcRepo = searchGrpc.NewSearchGrpcRepository(interviewConn)
	assert.NoError(suite.T(), err)

	suite.repo = mock.NewMockSearchRepository(suite.controller)
	suite.list, err = net.Listen("tcp", "127.0.0.1:8003")
	assert.NoError(suite.T(), err)
	suite.server = grpc.NewServer()
	searchRpc.RegisterSearchServer(suite.server, NewSearchServer(suite.repo))
}

func (suite *userSuite) TearDown() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestSearchPersons() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().SearchPersons(gomock.Any()).Times(1).Return([]*baseModels.Person{&baseModels.Person{}}, nil)

	_, err := suite.grpcRepo.SearchPersons(&baseModels.SearchParams{})

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchOrganizations() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().SearchOrganizations(gomock.Any()).Times(1).Return([]*baseModels.Organization{&baseModels.Organization{}}, nil)

	_, err := suite.grpcRepo.SearchOrganizations(&baseModels.SearchParams{})

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchVacancies() {
	go suite.server.Serve(suite.list)
	defer suite.server.Stop()

	suite.repo.EXPECT().SearchVacancies(gomock.Any()).Times(1).Return([]*baseModels.Vacancy{&baseModels.Vacancy{}}, nil)

	_, err := suite.grpcRepo.SearchVacancies(&baseModels.SearchParams{})

	assert.NoError(suite.T(), err)
}