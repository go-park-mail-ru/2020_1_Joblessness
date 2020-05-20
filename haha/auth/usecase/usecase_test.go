package authUseCase

//go:generate mockgen -destination=../repository/mock/auth.go -package=mock joblessness/haha/auth/interfaces AuthRepository

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/auth/repository/mock"
	"joblessness/haha/models/base"
	"testing"
)

type userSuite struct {
	suite.Suite
	controller   *gomock.Controller
	uc           *AuthUseCase
	person       baseModels.Person
	organization baseModels.Organization
	repo         *mock.MockAuthRepository
	sidEx        string
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
	defer suite.controller.Finish()

	suite.repo = mock.NewMockAuthRepository(suite.controller)
	suite.uc = NewAuthUseCase(suite.repo)

	suite.sidEx = "awdawd"

	suite.person = baseModels.Person{
		ID:        uint64(1),
		Login:     "user",
		Password:  "password",
		FirstName: "phone",
		Phone:     "name",
	}

	suite.organization = baseModels.Organization{
		Login:    "user",
		Password: "password",
		Name:     "name",
		Phone:    "phone",
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestAuthPersonFlow() {
	//RegisterPerson
	suite.repo.EXPECT().RegisterPerson(suite.person.Login, suite.person.Password, suite.person.FirstName).Return(nil).Times(1)
	suite.repo.EXPECT().DoesUserExists(suite.person.Login).Return(nil).Times(1)
	err := suite.uc.RegisterPerson(suite.person.Login, suite.person.Password, suite.person.FirstName)
	assert.NoError(suite.T(), err)

	//Login
	suite.repo.EXPECT().Login(suite.person.Login, suite.person.Password, gomock.Any()).Return(suite.person.ID, nil).Times(1)
	suite.repo.EXPECT().GetRole(suite.person.ID).Return("user", nil).Times(1)
	userId, _, sid, err := suite.uc.Login(suite.person.Login, suite.person.Password)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.person.ID, userId, "Id corrupted")
	assert.NotEmpty(suite.T(), sid, "No sid")

	//Logout
	suite.repo.EXPECT().Logout(suite.sidEx).Return(nil).Times(1)
	err = suite.uc.Logout(suite.sidEx)
	assert.NoError(suite.T(), err)

	//Check
	suite.repo.EXPECT().SessionExists(suite.sidEx).Return(suite.person.ID, nil).Times(1)
	userId, err = suite.uc.SessionExists(suite.sidEx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.person.ID, userId, "Id corrupted")
}

func (suite *userSuite) TestAuthOrganizationFlow() {
	suite.repo.EXPECT().RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name).Return(nil).Times(1)
	suite.repo.EXPECT().DoesUserExists(suite.organization.Login).Return(nil).Times(1)
	err := suite.uc.RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name)
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetRole() {
	suite.repo.EXPECT().GetRole(suite.organization.ID).
		Return(suite.organization.Name, nil).
		Times(1)

	_, err := suite.uc.GetRole(suite.organization.ID)
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestPersonSession() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(suite.person.ID, nil).
		Times(1)
	suite.repo.EXPECT().GetRole(suite.person.ID).
		Return("person", nil).
		Times(1)

	res, err := suite.uc.PersonSession(suite.sidEx)
	assert.Equal(suite.T(), suite.person.ID, res)
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestPersonSessionNoSession() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(uint64(0), errors.New("")).
		Times(1)

	_, err := suite.uc.PersonSession(suite.sidEx)
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestPersonSessionFailed() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(suite.person.ID, nil).
		Times(1)
	suite.repo.EXPECT().GetRole(suite.person.ID).
		Return("person", errors.New("")).
		Times(1)

	_, err := suite.uc.PersonSession(suite.sidEx)
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestPersonSessionWrongRole() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(suite.person.ID, nil).
		Times(1)
	suite.repo.EXPECT().GetRole(suite.person.ID).
		Return("organization", nil).
		Times(1)

	_, err := suite.uc.PersonSession(suite.sidEx)
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestOrgSession() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(suite.person.ID, nil).
		Times(1)
	suite.repo.EXPECT().GetRole(suite.person.ID).
		Return("organization", nil).
		Times(1)

	res, err := suite.uc.OrganizationSession(suite.sidEx)
	assert.Equal(suite.T(), suite.person.ID, res)
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestOrgSessionNoSession() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(uint64(0), errors.New("")).
		Times(1)

	_, err := suite.uc.OrganizationSession(suite.sidEx)
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestOrgSessionFailed() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(suite.person.ID, nil).
		Times(1)
	suite.repo.EXPECT().GetRole(suite.person.ID).
		Return("organization", errors.New("")).
		Times(1)

	_, err := suite.uc.OrganizationSession(suite.sidEx)
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestOrgSessionWrongRole() {
	suite.repo.EXPECT().SessionExists(suite.sidEx).
		Return(suite.person.ID, nil).
		Times(1)
	suite.repo.EXPECT().GetRole(suite.person.ID).
		Return("person", nil).
		Times(1)

	//RegisterOrganization
	suite.repo.EXPECT().RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name).Return(nil).Times(1)
	suite.repo.EXPECT().DoesUserExists(suite.organization.Login).Return(nil).Times(1)
	err := suite.uc.RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name)
	assert.NoError(suite.T(), err)
}
