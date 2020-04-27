package vacancyUseCase

import (
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	baseModels "joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
	mockRoom "joblessness/haha/utils/chat/mock"
	"joblessness/haha/vacancy/repository/mock"
	"testing"
)

//go:generate mockgen -destination=../repository/mock/vacancy.go -package=mock joblessness/haha/vacancy/interfaces VacancyRepository

type userSuite struct {
	suite.Suite
	controller   *gomock.Controller
	uc           *VacancyUseCase
	repo         *mock.MockVacancyRepository
	summary    baseModels.Summary
	education  pgModels.Education
	experience pgModels.Experience
	user       pgModels.User
	person     pgModels.Person
	vacancy   baseModels.Vacancy
	room *mockRoom.MockRoom
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
	defer suite.controller.Finish()

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
	suite.room = mockRoom.NewMockRoom(suite.controller)

	suite.repo = mock.NewMockVacancyRepository(suite.controller)
	suite.uc = NewVacancyUseCase(suite.repo, suite.room, bluemonday.UGCPolicy())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestAnnounceVacancy() {
	suite.repo.EXPECT().GetRelatedUsers(suite.vacancy.Organization.ID).Return([]uint64{1, 2}, suite.vacancy.Organization.Name, nil).Times(1)
	suite.room.EXPECT().SendGeneratedMessage(gomock.Any()).Times(2)

	err := suite.uc.announceVacancy(&suite.vacancy)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestCreateVacancy() {
	suite.repo.EXPECT().CreateVacancy(&suite.vacancy).Return(uint64(1), nil).Times(1)
	suite.repo.EXPECT().GetRelatedUsers(suite.vacancy.Organization.ID).Return([]uint64{1, 2}, suite.vacancy.Organization.Name, nil).Times(1)
	suite.room.EXPECT().SendGeneratedMessage(gomock.Any()).Times(2)

	_, err := suite.uc.CreateVacancy(&suite.vacancy)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetVacancies() {
	suite.repo.EXPECT().GetVacancies(1).Return(baseModels.Vacancies{&suite.vacancy}, nil).Times(1)

	_, err := suite.uc.GetVacancies("1")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetVacancy() {
	suite.repo.EXPECT().GetVacancy(suite.vacancy.ID).Return(&suite.vacancy, nil).Times(1)

	_, err := suite.uc.GetVacancy(suite.vacancy.ID)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestChangeVacancy() {
	suite.repo.EXPECT().ChangeVacancy(&suite.vacancy).Return(nil).Times(1)
	suite.repo.EXPECT().CheckAuthor(suite.vacancy.ID, suite.vacancy.Organization.ID).Return(nil).Times(1)

	err := suite.uc.ChangeVacancy(&suite.vacancy)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestDeleteVacancy() {
	suite.repo.EXPECT().DeleteVacancy(suite.vacancy.ID).Return(nil).Times(1)
	suite.repo.EXPECT().CheckAuthor(suite.vacancy.ID, suite.vacancy.Organization.ID).Return(nil).Times(1)

	err := suite.uc.DeleteVacancy(suite.vacancy.ID, suite.vacancy.Organization.ID)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetOrgVacancies() {
	suite.repo.EXPECT().GetOrgVacancies(uint64(1)).Return(baseModels.Vacancies{&suite.vacancy}, nil).Times(1)

	_, err := suite.uc.GetOrgVacancies(uint64(1))

	assert.NoError(suite.T(), err)
}