package vacancyRepoPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/models"
	"testing"
	"time"
)

type vacancySuite struct {
	suite.Suite
	rep *VacancyRepository
	db *sql.DB
	mock sqlmock.Sqlmock
	vacancy models.Vacancy
	user User
	organization Organization
}

func (suite *vacancySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewVacancyRepository(suite.db)

	suite.vacancy = models.Vacancy{
		ID:              3,
		Organization:     models.VacancyOrganization{
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

	suite.user = User{
		ID:             12,
		OrganizationID: 1,
		PersonID:       0,
		Tag:            "tag",
		Email:          "email",
		Phone:          "phone",
		Registered:     time.Time{},
		Avatar:         "avatar",
	}

	suite.organization = Organization{
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
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
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

	rows = sqlmock.NewRows([]string{"organization_id", "tag", "email", "phone", "avatar", "name", "site"}).
		AddRow(suite.user.OrganizationID, suite.user.Tag, suite.user.Email, suite.user.Phone, suite.user.Avatar,
			suite.organization.Name, suite.organization.Site)
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
		WithArgs(10, 9).
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
	assert.Equal(suite.T(), suite.vacancy, vacancy[0])
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

func (suite *vacancySuite) TestChangeVacancy() {
	suite.mock.
		ExpectExec("UPDATE vacancy SET organization_id =").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
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