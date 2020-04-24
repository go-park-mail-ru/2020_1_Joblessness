package authPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/status"
	authInterfaces "joblessness/haha/auth/interfaces"
	"joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
	"joblessness/haha/utils/salt"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	rep          *AuthRepository
	db           *sql.DB
	mock         sqlmock.Sqlmock
	person       baseModels.Person
	organization baseModels.Organization
}

func (suite *userSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewAuthRepository(suite.db)

	suite.person = baseModels.Person{
		ID:        1,
		Login:     "login",
		Password:  "password",
		FirstName: "first",
		LastName:  "name",
		Email:     "email",
		Phone:     "phone",
		Tag:       "tag",
	}

	suite.organization = baseModels.Organization{
		ID:       1,
		Login:    "login",
		Password: "password",
		Name:     "name",
		Site:     "site",
		Email:    "email",
		Phone:    "phone",
		Tag:      "tag",
	}
}

func (suite *userSuite) TearDown() {
	assert.NoError(suite.T(), suite.db.Close())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestDoesNotExists() {
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	suite.mock.
		ExpectQuery("SELECT count").
		WithArgs("name").
		WillReturnRows(rows)

	err := suite.rep.DoesUserExists("name")
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestDoesExists() {
	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	suite.mock.
		ExpectQuery("SELECT count").
		WithArgs("name").
		WillReturnRows(rows)

	err := suite.rep.DoesUserExists("name")
	e, ok := status.FromError(err)

	assert.True(suite.T(), ok)
	assert.True(suite.T(), e.Code() == authInterfaces.AlreadyExists)
}

func (suite *userSuite) TestDoesExistsErr() {
	suite.mock.
		ExpectQuery("SELECT count").
		WithArgs("name").
		WillReturnError(errors.New(""))

	err := suite.rep.DoesUserExists("name")
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestCreateUserNoId() {
	user, _ := pgModels.ToPgPerson(&suite.person)
	err := suite.rep.CreateUser(user.Login, user.Password, user.PersonID, 0)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestCreatePerson() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	suite.mock.
		ExpectQuery("INSERT INTO person ").
		WithArgs(suite.person.FirstName).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO users").
		WithArgs("login", sqlmock.AnyArg(), 0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.RegisterPerson(suite.person.Login, suite.person.Password, suite.person.FirstName)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestCreatePersonFailed() {
	suite.mock.
		ExpectQuery("INSERT INTO person ").
		WithArgs(suite.person.FirstName + " " + suite.person.LastName).
		WillReturnError(errors.New(""))

	err := suite.rep.RegisterPerson(suite.person.Login, suite.person.Password, suite.person.FirstName)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestCreateOrg() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	suite.mock.
		ExpectQuery("INSERT INTO organization").
		WithArgs(suite.organization.Name).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO users").
		WithArgs("login", sqlmock.AnyArg(), 1, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestCreateOrgFailed() {
	suite.mock.
		ExpectQuery("INSERT INTO organization").
		WithArgs(suite.organization.Name).
		WillReturnError(errors.New(""))

	err := suite.rep.RegisterOrganization(suite.organization.Login, suite.organization.Password, suite.organization.Name)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestLogin() {
	hashedPsw, _ := salt.HashAndSalt("password")
	rows := sqlmock.NewRows([]string{"id", "password"}).AddRow(1, hashedPsw)

	suite.mock.
		ExpectQuery("SELECT id, password").
		WithArgs(suite.person.Login).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO session").
		WithArgs(1, "sid", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.rep.Login(suite.person.Login, suite.person.Password, "sid")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestLoginWrongPassword() {
	rows := sqlmock.NewRows([]string{"id", "password"}).AddRow(1, "password")

	suite.mock.
		ExpectQuery("SELECT id, password").
		WithArgs(suite.person.Login).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO session").
		WithArgs(1, "sid", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.rep.Login(suite.person.Login, suite.person.Password, "sid")

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestLogout() {
	suite.mock.
		ExpectExec("DELETE FROM session").
		WithArgs("sid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.Logout("sid")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSessionExists() {
	rows := sqlmock.NewRows([]string{"id", "expires"}).AddRow(1, time.Now().Add(time.Hour))

	suite.mock.
		ExpectQuery("SELECT user_id, expires").
		WithArgs("sid").
		WillReturnRows(rows)

	userID, err := suite.rep.SessionExists("sid")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, uint64(1))
}

func (suite *userSuite) TestSessionExistsFailed() {
	suite.mock.
		ExpectQuery("SELECT user_id, expires").
		WithArgs("sid").
		WillReturnError(errors.New(""))

	_, err := suite.rep.SessionExists("sid")

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestSessionExistsExpired() {
	rows := sqlmock.NewRows([]string{"id", "expires"}).AddRow(1, time.Now().Add(-1*time.Hour))

	suite.mock.
		ExpectQuery("SELECT user_id, expires").
		WithArgs("sid").
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("DELETE FROM session").
		WithArgs("sid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.rep.SessionExists("sid")
	e, ok := status.FromError(err)

	assert.True(suite.T(), ok)
	assert.True(suite.T(), e.Code() == authInterfaces.WrongSID)
}
