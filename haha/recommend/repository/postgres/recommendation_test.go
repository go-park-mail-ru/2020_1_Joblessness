package recommendPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
}

func (suite *recommendSuite) SetupTest() {
	var err error
	suite.db, suite.sqlMock, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.vacancyRepository = vacancyPostgres.NewVacancyRepository(suite.db)
	suite.repository = NewRecommendRepository(suite.db, suite.vacancyRepository)
}

func (suite *recommendSuite) TearDown() {}

func TestSuite(t *testing.T) {
	suite.Run(t, new(recommendSuite))
}

func (suite *recommendSuite) TestNoUser() {
	_, _, err := suite.repository.GetRecommendedVacancies(1, 10, 0)
	assert.True(suite.T(), errors.Is(err, recommendInterfaces.ErrNoUser))
}
