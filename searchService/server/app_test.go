package searchServer

//go:generate mockgen -destination=../../haha/search/repository/mock/search.go -package=mock joblessness/haha/search/interfaces SearchRepository

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
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
	repo       *mock.MockSearchRepository
	server     *grpc.Server
	list       *bufconn.Listener
	conn       *grpc.ClientConn
}

func (suite *userSuite) bufDialer(context.Context, string) (net.Conn, error) {
	return suite.list.Dial()
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())

	suite.repo = mock.NewMockSearchRepository(suite.controller)
	buffer := 1024 * 1024
	suite.list = bufconn.Listen(buffer)
	suite.server = grpc.NewServer()
	searchRpc.RegisterSearchServer(suite.server, NewSearchServer(suite.repo))

	ctx := context.Background()
	var err error
	suite.conn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(suite.bufDialer), grpc.WithInsecure())

	suite.grpcRepo = searchGrpc.NewSearchGrpcRepository(suite.conn)
	assert.NoError(suite.T(), err)

	go func() {
		err = suite.server.Serve(suite.list)
	}()
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TearDown() {
	err := suite.conn.Close()
	assert.NoError(suite.T(), err)
	suite.conn.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestSearchPersons() {
	suite.repo.EXPECT().SearchPersons(gomock.Any()).Times(1).Return([]*baseModels.Person{&baseModels.Person{}}, nil)

	_, err := suite.grpcRepo.SearchPersons(&baseModels.SearchParams{})

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchOrganizations() {
	suite.repo.EXPECT().SearchOrganizations(gomock.Any()).Times(1).Return([]*baseModels.Organization{&baseModels.Organization{}}, nil)

	_, err := suite.grpcRepo.SearchOrganizations(&baseModels.SearchParams{})

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchVacancies() {
	suite.repo.EXPECT().SearchVacancies(gomock.Any()).Times(1).Return([]*baseModels.Vacancy{&baseModels.Vacancy{}}, nil)

	_, err := suite.grpcRepo.SearchVacancies(&baseModels.SearchParams{})

	assert.NoError(suite.T(), err)
}
