package authServer

//go:generate mockgen -destination=../../haha/auth/repository/mock/auth.go -package=mock joblessness/haha/auth/interfaces AuthRepository

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	authGrpc "joblessness/authService/grpc"
	"joblessness/haha/auth/repository/grpc"
	"joblessness/haha/auth/repository/mock"
	"net"
	"testing"
)

type userSuite struct {
	suite.Suite
	controller *gomock.Controller
	grpcRepo   *authGrpcRepository.AuthRepository
	repo       *mock.MockAuthRepository
	server     *grpc.Server
	list       *bufconn.Listener
	conn       *grpc.ClientConn
}

func (suite *userSuite) bufDialer(context.Context, string) (net.Conn, error) {
	return suite.list.Dial()
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())

	suite.repo = mock.NewMockAuthRepository(suite.controller)
	buffer := 1024 * 1024
	suite.list = bufconn.Listen(buffer)
	suite.server = grpc.NewServer()
	authGrpc.RegisterAuthServer(suite.server, NewAuthServer(suite.repo))

	ctx := context.Background()
	var err error
	suite.conn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(suite.bufDialer), grpc.WithInsecure())

	suite.grpcRepo = authGrpcRepository.NewRepository(suite.conn)
	assert.NoError(suite.T(), err)

	go func() {
		err = suite.server.Serve(suite.list)
	}()
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TearDown() {
	err := suite.conn.Close()
	assert.NoError(suite.T(), err)
	suite.server.Stop()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestRegisterPerson() {
	suite.repo.EXPECT().RegisterPerson("awd", "awda", "awda").Times(1).Return(nil)

	err := suite.grpcRepo.RegisterPerson("awd", "awda", "awda")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestRegisterOrganization() {
	suite.repo.EXPECT().RegisterOrganization("awd", "awda", "awda").Times(1).Return(nil)

	err := suite.grpcRepo.RegisterOrganization("awd", "awda", "awda")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestLogin() {
	suite.repo.EXPECT().Login("awd", "awda", "awda").Times(1).Return(uint64(2), nil)

	userID, err := suite.grpcRepo.Login("awd", "awda", "awda")

	assert.Equal(suite.T(), uint64(2), userID)
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestLogout() {
	suite.repo.EXPECT().Logout("awdaw").Times(1).Return(nil)

	err := suite.grpcRepo.Logout("awdaw")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSessionExists() {
	suite.repo.EXPECT().Logout("awdaw").Times(1).Return(nil)

	err := suite.grpcRepo.Logout("awdaw")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestDoesUserExists() {
	suite.repo.EXPECT().DoesUserExists("awdaw").Times(1).Return(nil)

	err := suite.grpcRepo.DoesUserExists("awdaw")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetRole() {
	suite.repo.EXPECT().GetRole(uint64(2)).Times(1).Return("role", nil)

	res, err := suite.grpcRepo.GetRole(uint64(2))

	assert.Equal(suite.T(), "role", res)
	assert.NoError(suite.T(), err)
}
