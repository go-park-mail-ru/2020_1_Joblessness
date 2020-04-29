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
	vacancy baseModels.Vacancy
}

func (suite *recommendSuite) SetupTest() {
	var err error
	suite.db, suite.sqlMock, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.vacancyRepository = vacancyPostgres.NewVacancyRepository(suite.db)
	suite.repository = NewRecommendRepository(suite.db, suite.vacancyRepository)

	suite.vacancy = baseModels.Vacancy{
		ID:               1,
		Organization:     baseModels.VacancyOrganization{
			ID:     1,
			Tag:    "tag",
			Email:  "email",
			Phone:  "phone",
			Avatar: "avatar",
			Name:   "name",
			Site:   "site",
		},
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

//func (suite *recommendSuite) TestGetRecommendedVacancies() {
//	rows := sqlmock.NewRows([]string{"id", "organization_id", "name", "description", "with_tax", "responsibilities", "conditions", "keywords", "salary_from", "salary_to"}).
//		AddRow(suite.vacancy.ID, suite.vacancy.Organization.ID, suite.vacancy.Name, suite.vacancy.Description, suite.vacancy.WithTax, suite.vacancy.Responsibilities, suite.vacancy.Conditions, suite.vacancy.Keywords, suite.vacancy.SalaryFrom, suite.vacancy.SalaryTo)
//	suite.sqlMock.
//		ExpectQuery("SELECT id, organization_id, name, description, with_tax, responsibilities, conditions, keywords, salary_from, salary_to").
//		WithArgs(pq.Array([]int{}), 10, 1).
//		WillReturnRows(rows)
//
//	_, _, err := suite.repository.GetRecommendedVacancies(1, 10, 0)
//	assert.NoError(suite.T(), err)
//}
