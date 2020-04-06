package authRepoPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/models"
	"joblessness/haha/utils/salt"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	rep *UserRepository
	db *sql.DB
	mock sqlmock.Sqlmock
	person models.Person
	organization models.Organization
}

func (suite *userSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewUserRepository(suite.db)

	suite.person = models.Person{
		ID: 1,
		Login:       "login",
		Password:    "password",
		FirstName:   "first",
		LastName:    "name",
		Email:       "email",
		Phone: "phone",
		Tag: "tag",
	}

	suite.organization = models.Organization{
		ID: 1,
		Login:       "login",
		Password:    "password",
		Name:   "name",
		Site:    "site",
		Email:       "email",
		Phone: "phone",
		Tag: "tag",
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
	assert.Equal(suite.T(), err, authInterfaces.ErrUserAlreadyExists)
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
	user, _ := toPostgresPerson(&suite.person)
	err := suite.rep.CreateUser(user)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestSaveAvatar() {
	suite.mock.
		ExpectExec("UPDATE users").
		WithArgs("avatar", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.SaveAvatarLink("avatar", 1)

	assert.NoError(suite.T(), err)
}


func (suite *userSuite) TestSaveAvatarEmpty() {
	err := suite.rep.SaveAvatarLink("", 1)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestCreatePerson() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	suite.mock.
		ExpectQuery("INSERT INTO person ").
		WithArgs(suite.person.FirstName + " " + suite.person.LastName, suite.person.Gender, suite.person.Birthday).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO users").
		WithArgs("login", sqlmock.AnyArg(), 0, 1, "email", "phone", "tag").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.CreatePerson(&suite.person)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestCreatePersonFailed() {
	suite.mock.
		ExpectQuery("INSERT INTO person ").
		WithArgs(suite.person.FirstName + " " + suite.person.LastName).
		WillReturnError(errors.New(""))

	err := suite.rep.CreatePerson(&suite.person)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestCreateOrg() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	suite.mock.
		ExpectQuery("INSERT INTO organization").
		WithArgs(suite.organization.Name, suite.organization.Site, suite.organization.About).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO users").
		WithArgs("login", sqlmock.AnyArg(), 1, 0, "email", "phone", "tag").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.CreateOrganization(&suite.organization)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestCreateOrgFailed() {
	suite.mock.
		ExpectQuery("INSERT INTO organization").
		WithArgs(suite.organization.Name).
		WillReturnError(errors.New(""))

	err := suite.rep.CreateOrganization(&suite.organization)

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
	rows := sqlmock.NewRows([]string{"id", "expires"}).AddRow(1, time.Now().Add(-1 * time.Hour))

	suite.mock.
		ExpectQuery("SELECT user_id, expires").
		WithArgs("sid").
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("DELETE FROM session").
		WithArgs("sid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.rep.SessionExists("sid")

	assert.Equal(suite.T(), err, authInterfaces.ErrWrongSID)
}

func (suite *userSuite) TestGetPerson() {
	rows := sqlmock.NewRows([]string{"login", "person_id", "email", "phone", "avatar", "tag"})
	rows = rows.AddRow(suite.person.Login, suite.person.ID,
		suite.person.Email, suite.person.Phone, suite.person.Avatar, suite.person.Tag)
	suite.mock.
		ExpectQuery("SELECT login").
		WithArgs(1).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"name", "gender", "birthday"})
	rows = rows.AddRow(suite.person.LastName + " " + suite.person.LastName, suite.person.Gender, suite.person.Birthday)
	suite.mock.
		ExpectQuery("SELECT name").
		WithArgs(suite.person.ID).
		WillReturnRows(rows)

	_, err := suite.rep.GetPerson(uint64(1))

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetPersonFailedOne() {
	rows := sqlmock.NewRows([]string{"login", "person_id", "email", "phone", "avatar"})
	rows = rows.AddRow(suite.person.Login, suite.person.ID,
		suite.person.Email, suite.person.Phone, suite.person.Avatar)
	suite.mock.
		ExpectQuery("SELECT login, person_id, email, phone, avatar").
		WithArgs(12).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetPerson(uint64(12))

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestGetPersonFailedTwo() {
	rows := sqlmock.NewRows([]string{"login", "person_id", "email", "phone", "avatar"})
	rows = rows.AddRow(suite.person.Login, suite.person.ID,
		suite.person.Email, suite.person.Phone, suite.person.Avatar)
	suite.mock.
		ExpectQuery("SELECT login, person_id, email, phone, avatar").
		WithArgs(12).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"name"})
	rows = rows.AddRow(suite.person.LastName + " " + suite.person.LastName)
	suite.mock.
		ExpectQuery("SELECT name").
		WithArgs(suite.person.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetPerson(uint64(12))

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestChangePerson() {
	rows := sqlmock.NewRows([]string{"person_id"})
	rows = rows.AddRow(suite.person.ID)
	suite.mock.
		ExpectQuery("SELECT person_id").
		WithArgs(suite.person.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectExec("UPDATE person").
		WithArgs(suite.person.FirstName + " " + suite.person.LastName, suite.person.Gender, nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec("UPDATE user").
		WithArgs(suite.person.Password, suite.person.Tag, suite.person.Email, suite.person.Phone, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.ChangePerson(suite.person)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestChangePersonFailedOne() {
	rows := sqlmock.NewRows([]string{"person_id"})
	rows = rows.AddRow(suite.person.ID)
	suite.mock.
		ExpectQuery("SELECT person_id").
		WithArgs(suite.person.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangePerson(suite.person)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestChangePersonFailedTwo() {
	rows := sqlmock.NewRows([]string{"person_id"})
	rows = rows.AddRow(suite.person.ID)
	suite.mock.
		ExpectQuery("SELECT person_id").
		WithArgs(suite.person.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectExec("UPDATE person").
		WithArgs(suite.person.FirstName + " " + suite.person.LastName, 12).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangePerson(suite.person)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestGetOrganization() {
	rows := sqlmock.NewRows([]string{"login", "organization_id", "email", "phone", "avatar", "tag"})
	rows = rows.AddRow(suite.organization.Login, suite.organization.ID,
		suite.organization.Email, suite.organization.Phone, suite.organization.Avatar, suite.organization.Tag)
	suite.mock.
		ExpectQuery("SELECT login,").
		WithArgs(12).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"name", "site", "about"})
	rows = rows.AddRow(suite.organization.Name, suite.organization.Site, suite.organization.About)
	suite.mock.
		ExpectQuery("SELECT name, site, about").
		WithArgs(suite.organization.ID).
		WillReturnRows(rows)

	result, err := suite.rep.GetOrganization(uint64(12))

	assert.NoError(suite.T(), err)
	assert.ObjectsAreEqual(result, suite.organization)
}

func (suite *userSuite) TestGetOrganizationFailedOne() {
	rows := sqlmock.NewRows([]string{"login", "organization_id", "email", "phone", "avatar"})
	rows = rows.AddRow(suite.organization.Login, suite.organization.ID,
		suite.organization.Email, suite.organization.Phone, suite.organization.Avatar)
	suite.mock.
		ExpectQuery("SELECT login, organization_id, email, phone, avatar").
		WithArgs(12).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetOrganization(uint64(12))

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestGetOrganizationFailedTwo() {
	rows := sqlmock.NewRows([]string{"login", "password", "organization_id", "email", "phone", "avatar"})
	rows = rows.AddRow(suite.organization.Login, suite.organization.Password, suite.organization.ID,
		suite.organization.Email, suite.organization.Phone, suite.organization.Avatar)
	suite.mock.
		ExpectQuery("SELECT login, password, organization_id, email, phone, avatar").
		WithArgs(12).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT name").
		WithArgs(suite.organization.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetOrganization(uint64(12))

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestChangeOrganization() {
	rows := sqlmock.NewRows([]string{"organization_id"})
	rows = rows.AddRow(suite.organization.ID)
	suite.mock.
		ExpectQuery("SELECT organization_id").
		WithArgs(suite.organization.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectExec("UPDATE organization").
		WithArgs(suite.organization.Name, suite.organization.Site, suite.organization.About, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec("UPDATE user").
		WithArgs(suite.organization.Password, suite.organization.Tag, suite.organization.Email, suite.organization.Phone, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.ChangeOrganization(suite.organization)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestChangeOrganizationFailedOne() {
	rows := sqlmock.NewRows([]string{"organization_id"})
	rows = rows.AddRow(suite.organization.ID)
	suite.mock.
		ExpectQuery("SELECT organization_id").
		WithArgs(suite.organization.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangeOrganization(suite.organization)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestChangeOrganizationFailedTwo() {
	rows := sqlmock.NewRows([]string{"organization_id"})
	rows = rows.AddRow(suite.organization.ID)
	suite.mock.
		ExpectQuery("SELECT organization_id").
		WithArgs(suite.organization.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectExec("UPDATE organization").
		WithArgs(suite.organization.Name, 12).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangeOrganization(suite.organization)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestGetOrgList() {
	rows := sqlmock.NewRows([]string{"userId", "name", "site"})
	for i := 1; i < 5; i++ {
		rows = rows.AddRow(uint64(i), suite.organization.Name, suite.organization.Site)
	}

	suite.mock.
		ExpectQuery("SELECT users.id as userId, name, site").
		WithArgs(10, 10).
		WillReturnRows(rows)

	result, err := suite.rep.GetListOfOrgs(1)

	assert.Equal(suite.T(), 4, len(result))
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetOrgListFailed() {
	rows := sqlmock.NewRows([]string{"userId", "name", "site"})
	for i := 1; i < 5; i++ {
		rows = rows.AddRow(uint64(i), suite.organization.Name, suite.organization.Site)
	}

	suite.mock.
		ExpectQuery("SELECT users.id as userId, name, site").
		WithArgs(0, 9).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetListOfOrgs(1)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestLikeUserLike() {
	suite.mock.
		ExpectExec("INSERT INTO favorite").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := suite.rep.SetOrDeleteLike(suite.person.ID, suite.person.ID)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), res)
}

func (suite *userSuite) TestLikeUserLikeFailed() {
	suite.mock.
		ExpectExec("INSERT INTO favorite").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.SetOrDeleteLike(suite.person.ID, suite.person.ID)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestLikeExists() {
	rows := sqlmock.NewRows([]string{"count"}).
	AddRow(1)
	suite.mock.
		ExpectQuery("SELECT count").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnRows(rows)

	res, err := suite.rep.LikeExists(suite.person.ID, suite.person.ID)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), res)
}

func (suite *userSuite) TestLikeExistsFailed() {
	suite.mock.
		ExpectExec("SELECT count").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.LikeExists(suite.person.ID, suite.person.ID)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestLikeUserDis() {
	suite.mock.
		ExpectExec("INSERT INTO favorite").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))
	suite.mock.
		ExpectExec("DELETE FROM favorite").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := suite.rep.SetOrDeleteLike(suite.person.ID, suite.person.ID)

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), res)
}

func (suite *userSuite) TestLikeUserDisFailed() {
	suite.mock.
		ExpectExec("INSERT INTO favorite").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))
	suite.mock.
		ExpectExec("DELETE FROM favorite").
		WithArgs(suite.person.ID, suite.person.ID).
		WillReturnError(errors.New(""))
	_, err := suite.rep.SetOrDeleteLike(suite.person.ID, suite.person.ID)

	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestFavoritePer() {
	rows := sqlmock.NewRows([]string{"userId", "tag", "person_id"})
	rows = rows.AddRow(uint64(1), "tag", uint64(2))
	suite.mock.
		ExpectQuery("SELECT u.id, u.tag, u.person_id").
		WithArgs(suite.person.ID).
		WillReturnRows(rows)

	res, err := suite.rep.GetUserFavorite(suite.person.ID)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), res[0].IsPerson)
}

func (suite *userSuite) TestFavoriteOrg() {
	rows := sqlmock.NewRows([]string{"userId", "tag", "person_id"})
	rows = rows.AddRow(uint64(1), "tag", nil)
	suite.mock.
		ExpectQuery("SELECT u.id, u.tag, u.person_id").
		WithArgs(suite.person.ID).
		WillReturnRows(rows)

	res, err := suite.rep.GetUserFavorite(suite.person.ID)

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), res[0].IsPerson)
}

func (suite *userSuite) TestFavoriteFailed() {
	suite.mock.
		ExpectQuery("SELECT u.id, u.tag, u.person_id").
		WithArgs(suite.person.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetUserFavorite(suite.person.ID)

	assert.Error(suite.T(), err)
}