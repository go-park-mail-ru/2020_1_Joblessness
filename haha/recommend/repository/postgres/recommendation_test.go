package recommendPostgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	baseModels "joblessness/haha/models/base"
	"joblessness/haha/recommend/interfaces"
	"joblessness/haha/vacancy/interfaces"
	"joblessness/haha/vacancy/repository/postgres"
	"testing"
)

type recommendSuite struct {
	suite.Suite
	repository        recommendInterfaces.RecommendRepository
	vacancyRepository vacancyInterfaces.VacancyRepository
	db                *sql.DB
	sqlMock           sqlmock.Sqlmock
	vacancy           baseModels.Vacancy
	organization      baseModels.VacancyOrganization
}

func (suite *recommendSuite) SetupTest() {
	var err error
	suite.db, suite.sqlMock, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.vacancyRepository = vacancyPostgres.NewVacancyRepository(suite.db)
	suite.repository = NewRecommendRepository(suite.db, suite.vacancyRepository)

	suite.organization = baseModels.VacancyOrganization{
		ID:     1,
		Tag:    "tag",
		Email:  "email",
		Phone:  "phone",
		Avatar: "avatar",
		Name:   "name",
		Site:   "site",
	}

	suite.vacancy = baseModels.Vacancy{
		ID: 1,
		Name:             "name",
		Description:      "description",
		SalaryFrom:       10000,
		SalaryTo:         20000,
		WithTax:          false,
		Responsibilities: "responsibilities",
		Conditions:       "conditions",
		Keywords:         "keywords",
	}
}

func (suite *recommendSuite) TearDown() {}

func TestSuite(t *testing.T) {
	suite.Run(t, new(recommendSuite))
}

func (suite *recommendSuite) TestGetPopularVacancies() {
	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "with_tax", "responsibilities", "conditions", "keywords", "salary_from", "salary_to", "count"}).
		AddRow(suite.vacancy.ID, suite.organization.ID, suite.vacancy.Name, suite.vacancy.Description, suite.vacancy.WithTax, suite.vacancy.Responsibilities, suite.vacancy.Conditions, suite.vacancy.Keywords, suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo, 1)
	suite.sqlMock.
		ExpectQuery("SELECT v.id, v.organization_id, v.name, v.description, v.with_tax, v.responsibilities, v.conditions, v.keywords, v.salary_from, v.salary_to, COUNT\\(\\*\\) count").
		WithArgs(10, 0).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"organization_id", "tag", "email", "phone", "avatar", "name", "site"}).
		AddRow(suite.organization.ID, suite.organization.Tag, suite.organization.Email, suite.organization.Phone, suite.organization.Avatar, suite.organization.Name, suite.organization.Site)
	suite.sqlMock.
		ExpectQuery("SELECT organization_id, tag, email, phone, avatar, name, site").
		WithArgs(1).
		WillReturnRows(rows)

	_, err := suite.repository.GetPopularVacancies(10, 0)
	assert.NoError(suite.T(), err)
}
