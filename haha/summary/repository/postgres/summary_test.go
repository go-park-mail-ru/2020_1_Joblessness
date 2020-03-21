package summaryRepoPostgres

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

type summarySuite struct {
	suite.Suite
	rep *SummaryRepository
	db *sql.DB
	mock sqlmock.Sqlmock
	summary models.Summary
	education Education
	experience Experience
	user User
	person Person
}

func (suite *summarySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewSummaryRepository(suite.db)

	suite.summary = models.Summary{
		ID:          3,
		Author:      models.Author{
			ID:        12,
			Tag:       "tag",
			Email:     "email",
			Phone:     "phone",
			Avatar:    "avatar",
			FirstName: "first",
			LastName:  "name",
			Gender:    "gender",
			Birthday:  "birth",
		},
		Keywords:    "key",
		Educations:  []models.Education{
			models.Education{
				Institution: "was",
				Speciality:  "is",
				Graduated:   "yes",
				Type:        "none",
			},
		},
		Experiences: []models.Experience{
			models.Experience{
				CompanyName:      "comp",
				Role:             "role",
				Responsibilities: "response",
				Start:            "start",
				Stop:             "stop",
			},
		},
	}

	suite.user = User{
		ID:             12,
		OrganizationID: 0,
		PersonID:       5,
		Tag:            sql.NullString{},
		Email:          sql.NullString{},
		Phone:          sql.NullString{},
		Registered:     sql.NullTime{},
		Avatar:         sql.NullString{},
	}

	suite.person = Person{
		ID:       sql.NullString{},
		Name:     sql.NullString{},
		Gender:   sql.NullString{},
		Birthday: sql.NullTime{},
	}
}

func (suite *summarySuite) TearDown() {
	assert.NoError(suite.T(), suite.db.Close())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(summarySuite))
}

func (suite *summarySuite) TestCreateVacancy() {
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

func (suite *summarySuite) TestCreateVacancyFailed() {
	suite.mock.
		ExpectQuery("INSERT INTO vacancy").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax,
			suite.vacancy.Responsibilities, suite.vacancy.Conditions, suite.vacancy.Keywords).
		WillReturnError(errors.New(""))

	_, err := suite.rep.CreateVacancy(&suite.vacancy)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetVacancy() {
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

func (suite *summarySuite) TestGetVacancyFailedOne() {
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WithArgs(suite.vacancy.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetVacancy(suite.vacancy.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetVacancyFailedTwo() {
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

func (suite *summarySuite) TestGetVacancies() {
	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "salary_from", "salary_to", "with_tax",
		"responsibilities", "conditions", "keywords"}).
		AddRow(suite.vacancy.ID, suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords)
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"organization_id", "tag", "email", "phone", "avatar", "name", "site"}).
		AddRow(suite.user.OrganizationID, suite.user.Tag, suite.user.Email, suite.user.Phone, suite.user.Avatar,
			suite.organization.Name, suite.organization.Site)
	suite.mock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(suite.user.ID).
		WillReturnRows(rows)

	vacancy, err := suite.rep.GetVacancies()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.vacancy, vacancy[0])
}

func (suite *summarySuite) TestGetVacanciesFailedOne() {
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WillReturnError(errors.New(""))

	rows := sqlmock.NewRows([]string{"organization_id", "tag", "email", "phone", "avatar", "name", "site"}).
		AddRow(suite.user.OrganizationID, suite.user.Tag, suite.user.Email, suite.user.Phone, suite.user.Avatar,
			suite.organization.Name, suite.organization.Site)
	suite.mock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(suite.user.ID).
		WillReturnRows(rows)

	_, err := suite.rep.GetVacancies()
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetVacanciesFailedTwo() {
	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "salary_from", "salary_to", "with_tax",
		"responsibilities", "conditions", "keywords"}).
		AddRow(suite.vacancy.ID, suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords)
	suite.mock.
		ExpectQuery("SELECT id, organization_id, name, description, salary_from, salary_to, with_tax").
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(suite.user.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetVacancies()
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestChangeVacancy() {
	suite.mock.
		ExpectExec("UPDATE vacancy SET organization_id =").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords, suite.vacancy.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.ChangeVacancy(&suite.vacancy)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestChangeVacancyFailed() {
	suite.mock.
		ExpectExec("UPDATE vacancy SET organization_id =").
		WithArgs(suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description,
			suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax, suite.vacancy.Responsibilities,
			suite.vacancy.Conditions, suite.vacancy.Keywords, suite.vacancy.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangeVacancy(&suite.vacancy)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestDeleteVacancy() {
	suite.mock.
		ExpectExec("DELETE FROM vacancy").
		WithArgs(suite.vacancy.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.DeleteVacancy(suite.vacancy.ID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestDeleteVacancyFailed() {
	suite.mock.
		ExpectExec("DELETE FROM vacancy").
		WithArgs(suite.vacancy.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.DeleteVacancy(suite.vacancy.ID)
	assert.Error(suite.T(), err)
}