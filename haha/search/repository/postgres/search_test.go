package searchPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/models/base"
	"testing"
)

type userSuite struct {
	suite.Suite
	rep          *SearchRepository
	db           *sql.DB
	mock         sqlmock.Sqlmock
	person       baseModels.Person
	organization baseModels.Organization
	vacancy      baseModels.Vacancy
}

func (suite *userSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewSearchRepository(suite.db)

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

func (suite *userSuite) TestSearchPersons() {
	rows := sqlmock.NewRows([]string{"userId", "name", "tag", "avatar"}).
		AddRow(suite.person.ID, suite.person.FirstName+" "+suite.person.LastName, suite.person.Tag,
			suite.person.Avatar)
	suite.mock.
		ExpectQuery("SELECT users.id as userId, p.name, tag, avatar").
		WithArgs("req", 10, 10).
		WillReturnRows(rows)

	_, err := suite.rep.SearchPersons("req", "1", "true")
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchPersonsFailed() {
	suite.mock.
		ExpectQuery("SELECT users.id as userId, p.name, tag, avatar").
		WithArgs("req", 10, 10).
		WillReturnError(errors.New(""))

	_, err := suite.rep.SearchPersons("req", "1", "true")
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestSearchOrganization() {
	rows := sqlmock.NewRows([]string{"userId", "name", "tag", "avatar"}).
		AddRow(suite.organization.ID, suite.organization.Name, suite.organization.Tag, suite.organization.Avatar)
	suite.mock.
		ExpectQuery("SELECT users.id as userId, name, tag, avatar").
		WithArgs("req", 10, 10).
		WillReturnRows(rows)

	_, err := suite.rep.SearchOrganizations("req", "1", "true")
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchOrganizationFailed() {
	suite.mock.
		ExpectQuery("SELECT users.id as userId, name, tag, avatar").
		WithArgs("req", 10, 10).
		WillReturnError(errors.New(""))

	_, err := suite.rep.SearchOrganizations("req", "1", "true")
	assert.Error(suite.T(), err)
}

func (suite *userSuite) TestSearchVacancy() {
	rows := sqlmock.NewRows([]string{"id", "name", "id", "name", "keywords", "salary_from", "salary_to", "with_tax"}).
		AddRow(suite.organization.ID, suite.organization.Name, suite.vacancy.ID, suite.vacancy.Name,
			suite.vacancy.Keywords, suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, suite.vacancy.WithTax)
	suite.mock.
		ExpectQuery("SELECT users.id, o.name, v.id, v.name, v.keywords, v.salary_from, v.salary_to, v.with_tax").
		WithArgs("req", 10, 10).
		WillReturnRows(rows)

	_, err := suite.rep.SearchVacancies("req", "1", "true")
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchVacancyFailed() {
	suite.mock.
		ExpectQuery("SELECT users.id, o.name, v.id, v.name, v.keywords, v.salary_from, v.salary_to, v.with_tax").
		WithArgs("req", 10, 10).
		WillReturnError(errors.New(""))

	_, err := suite.rep.SearchVacancies("req", "1", "true")
	assert.Error(suite.T(), err)
}
