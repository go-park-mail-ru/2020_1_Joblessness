package searchUseCase

import (
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	baseModels "joblessness/haha/models/base"
	"joblessness/haha/search/repository/mock"
	"testing"
)

type userSuite struct {
	suite.Suite
	controller   *gomock.Controller
	uc           *SearchUseCase
	person       baseModels.Person
	organization baseModels.Organization
	repo         *mock.MockSearchRepository
	sidEx        string
	params       baseModels.SearchParams
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
	defer suite.controller.Finish()

	suite.repo = mock.NewMockSearchRepository(suite.controller)
	suite.uc = NewSearchUseCase(suite.repo, bluemonday.UGCPolicy())

	suite.params = baseModels.SearchParams{
		Request: "",
		Since:   "",
		Desc:    "",
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestSearchPerson() {
	suite.repo.EXPECT().SearchPersons(&suite.params).Return([]*baseModels.Person{}, nil).Times(1)

	_, err := suite.uc.Search("person", "", "", "")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchOrg() {
	suite.repo.EXPECT().SearchOrganizations(&suite.params).Return([]*baseModels.Organization{}, nil).Times(1)

	_, err := suite.uc.Search("organization", "", "", "")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchVacancy() {
	suite.repo.EXPECT().SearchVacancies(&suite.params).Return([]*baseModels.Vacancy{}, nil).Times(1)

	_, err := suite.uc.Search("vacancy", "", "", "")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSearchAll() {
	suite.repo.EXPECT().SearchVacancies(&suite.params).Return([]*baseModels.Vacancy{}, nil).Times(1)
	suite.repo.EXPECT().SearchOrganizations(&suite.params).Return([]*baseModels.Organization{}, nil).Times(1)
	suite.repo.EXPECT().SearchPersons(&suite.params).Return([]*baseModels.Person{}, nil).Times(1)

	_, err := suite.uc.Search("", "", "", "")

	assert.NoError(suite.T(), err)
}
