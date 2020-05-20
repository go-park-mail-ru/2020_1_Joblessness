package vacancyPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
	vacancyInterfaces "joblessness/haha/vacancy/interfaces"
	"testing"
	"time"
)

type vacancySuite struct {
	suite.Suite
	rep          *VacancyRepository
	db           *sql.DB
	mock         sqlmock.Sqlmock
	vacancy      baseModels.Vacancy
	user         pgModels.User
	organization pgModels.Organization
}

func (suite *vacancySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewVacancyRepository(suite.db)

	suite.vacancy = baseModels.Vacancy{
		ID: 3,
		Organization: baseModels.VacancyOrganization{
			ID:     12,
			Tag:    "tag",
			Email:  "email",
			Phone:  "phone",
			Avatar: "avatar",
			Name:   "name",
			Site:   "site",
		},
		Name:             "vacancy",
		Description:      "description",
		SalaryFrom:       50,
		SalaryTo:         100,
		WithTax:          false,
		Responsibilities: "all",
		Conditions:       "some",
		Keywords:         "word",
	}

	suite.user = pgModels.User{
		ID:             12,
		OrganizationID: 1,
		PersonID:       0,
		Tag:            "tag",
		Email:          "email",
		Phone:          "phone",
		Registered:     time.Time{},
		Avatar:         "avatar",
	}

	suite.organization = pgModels.Organization{
		ID:   1,
		Name: "name",
		Site: "site",
	}
}

func (suite *vacancySuite) TearDown() {
	assert.NoError(suite.T(), suite.db.Close())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(vacancySuite))
}

func (suite *vacancySuite) TestGetRelatedUsers() {
	rows := sqlmock.NewRows([]string{"name"}).AddRow(suite.organization.Name)
	suite.mock.
		ExpectQuery("SELECT o.name").
		WithArgs(suite.vacancy.Organization.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"user_id"}).AddRow(uint64(1)).AddRow(uint64(2))
	suite.mock.
		ExpectQuery("SELECT f.user_id").
		WithArgs(suite.vacancy.Organization.ID).
		WillReturnRows(rows)

	_, name, err := suite.rep.GetRelatedUsers(suite.vacancy.Organization.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.organization.Name, name)
}

func (suite *vacancySuite) TestGetRelatedUsersFailedOne() {
	suite.mock.
		ExpectQuery("SELECT o.name").
		WithArgs(suite.vacancy.Organization.ID).
		WillReturnError(errors.New(""))

	_, _, err := suite.rep.GetRelatedUsers(suite.vacancy.Organization.ID)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestGetRelatedUsersFailedTwo() {
	rows := sqlmock.NewRows([]string{"name"}).AddRow(suite.organization.Name)
	suite.mock.
		ExpectQuery("SELECT o.name").
		WithArgs(suite.vacancy.Organization.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT f.user_id").
		WithArgs(suite.vacancy.Organization.ID).
		WillReturnError(errors.New(""))

	_, _, err := suite.rep.GetRelatedUsers(suite.vacancy.Organization.ID)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestCreateVacancy() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(uint64(3))

	suite.mock.
		ExpectQuery("INSERT INTO vacancy").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax,
			suite.vacancy.Responsibilities, suite.vacancy.Conditions, suite.vacancy.Keywords).
		WillReturnRows(rows)

	vacancyID, err := suite.rep.CreateVacancy(&suite.vacancy)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.vacancy.ID, vacancyID)
}

func (suite *vacancySuite) TestCreateVacancyFailed() {
	suite.mock.
		ExpectQuery("INSERT INTO vacancy").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax,
			suite.vacancy.Responsibilities, suite.vacancy.Conditions, suite.vacancy.Keywords).
		WillReturnError(errors.New(""))

	_, err := suite.rep.CreateVacancy(&suite.vacancy)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestGetVacancy() {
	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "salary_from", "salary_to", "with_tax",
		"responsibilities", "conditions", "keywords"}).
		AddRow(suite.vacancy.ID, suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords)
	suite.mock.
		ExpectQuery("SELECT v.id, v.organization_id, v.name, v.description, v.salary_from").
		WithArgs(suite.vacancy.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"organization_id", "tag", "email", "phone", "avatar", "name", "site"}).
		AddRow(suite.user.OrganizationID, suite.user.Tag, suite.user.Email, suite.user.Phone, suite.user.Avatar,
			suite.organization.Name, suite.organization.Site)
	suite.mock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(suite.user.ID).
		WillReturnRows(rows)

	vacancy, err := suite.rep.GetVacancy(suite.vacancy.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.vacancy, *vacancy)
}

func (suite *vacancySuite) TestGetVacancyFailedOne() {
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WithArgs(suite.vacancy.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetVacancy(suite.vacancy.ID)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestGetVacancyFailedTwo() {
	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "salary_from", "salary_to", "with_tax",
		"responsibilities", "conditions", "keywords"}).
		AddRow(suite.vacancy.ID, suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords)
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WithArgs(suite.vacancy.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(suite.user.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetVacancy(suite.vacancy.ID)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestGetVacancies() {
	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "salary_from", "salary_to", "with_tax",
		"responsibilities", "conditions", "keywords"}).
		AddRow(suite.vacancy.ID, suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords)
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WithArgs(10, 10).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"organization_id", "tag", "email", "phone", "avatar", "name", "site"}).
		AddRow(suite.user.OrganizationID, suite.user.Tag, suite.user.Email, suite.user.Phone, suite.user.Avatar,
			suite.organization.Name, suite.organization.Site)
	suite.mock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(suite.user.ID).
		WillReturnRows(rows)

	vacancy, err := suite.rep.GetVacancies(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), &suite.vacancy, vacancy[0])
}

func (suite *vacancySuite) TestGetVacanciesFailedOne() {
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetVacancies(1)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestGetVacanciesFailedTwo() {
	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "salary_from", "salary_to", "with_tax",
		"responsibilities", "conditions", "keywords"}).
		AddRow(suite.vacancy.ID, suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords)
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WithArgs(10, 9).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(suite.user.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetVacancies(1)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestCheckAuthor() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(true)
	suite.mock.
		ExpectQuery("SELECT organization_id").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.ID).
		WillReturnRows(rows)

	err := suite.rep.CheckAuthor(suite.vacancy.ID, suite.vacancy.Organization.ID)
	assert.NoError(suite.T(), err)
}

func (suite *vacancySuite) TestCheckAuthorFailed() {
	suite.mock.
		ExpectQuery("SELECT organization_id").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.CheckAuthor(suite.vacancy.ID, suite.vacancy.Organization.ID)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestCheckAuthorNotAuthor() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(false)
	suite.mock.
		ExpectQuery("SELECT organization_id").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.ID).
		WillReturnRows(rows)

	err := suite.rep.CheckAuthor(suite.vacancy.ID, suite.vacancy.Organization.ID)
	assert.EqualError(suite.T(), vacancyInterfaces.ErrOrgIsNotOwner, err.Error())
}

func (suite *vacancySuite) TestChangeVacancy() {
	suite.mock.
		ExpectExec("UPDATE vacancy SET name =").
		WithArgs(suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords, suite.vacancy.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.ChangeVacancy(&suite.vacancy)
	assert.NoError(suite.T(), err)
}

func (suite *vacancySuite) TestChangeVacancyFailed() {
	suite.mock.
		ExpectExec("UPDATE vacancy SET organization_id =").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords, suite.vacancy.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangeVacancy(&suite.vacancy)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestDeleteVacancy() {
	suite.mock.
		ExpectExec("DELETE FROM vacancy").
		WithArgs(suite.vacancy.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.DeleteVacancy(suite.vacancy.ID)
	assert.NoError(suite.T(), err)
}

func (suite *vacancySuite) TestDeleteVacancyFailed() {
	suite.mock.
		ExpectExec("DELETE FROM vacancy").
		WithArgs(suite.vacancy.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.DeleteVacancy(suite.vacancy.ID)
	assert.Error(suite.T(), err)
}

func (suite *vacancySuite) TestGetOrgVacancies() {
	rows := sqlmock.NewRows([]string{"id", "name", "salary_from", "salary_to", "with_tax", "keywords"}).
		AddRow(suite.vacancy.ID, suite.vacancy.Name, suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo,
			suite.vacancy.WithTax, suite.vacancy.Keywords)
	suite.mock.
		ExpectQuery("SELECT id, name, salary_from, salary_to, with_tax").
		WithArgs(1).
		WillReturnRows(rows)

	_, err := suite.rep.GetOrgVacancies(1)
	assert.NoError(suite.T(), err)
}

func (suite *vacancySuite) TestGetOrgVacanciesFailedOne() {
	suite.mock.
		ExpectQuery("SELECT id, name, salary_from, salary_to, with_tax").
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetVacancies(1)
	assert.Error(suite.T(), err)
}
