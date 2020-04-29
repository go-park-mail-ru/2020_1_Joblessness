package summaryUseCase

//go:generate mockgen -destination=../repository/mock/summary.go -package=mock joblessness/haha/summary/interfaces SummaryRepository

import (
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	baseModels "joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
	summaryInterfaces "joblessness/haha/summary/interfaces"
	"joblessness/haha/summary/repository/mock"
	"testing"
	"time"
)

type userSuite struct {
	suite.Suite
	controller *gomock.Controller
	uc         *SummaryUseCase
	repo       *mock.MockSummaryRepository
	summary    baseModels.Summary
	education  pgModels.Education
	experience pgModels.Experience
	user       pgModels.User
	person     pgModels.Person
	response   baseModels.VacancyResponse
	sendSum    baseModels.SendSummary
}

func (suite *userSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
	defer suite.controller.Finish()

	suite.summary = baseModels.Summary{
		ID: 3,
		Author: baseModels.Author{
			ID:        12,
			Tag:       "tag",
			Email:     "email",
			Phone:     "phone",
			Avatar:    "avatar",
			FirstName: "name",
			LastName:  "first",
			Gender:    "gender",
		},
		Keywords: "key",
		Educations: []baseModels.Education{
			baseModels.Education{
				Institution: "was",
				Speciality:  "is",
				Type:        "none",
			},
		},
		Experiences: []baseModels.Experience{
			baseModels.Experience{
				CompanyName:      "comp",
				Role:             "role",
				Responsibilities: "response",
				Start:            time.Now(),
				Stop:             time.Now().AddDate(1, 1, 1),
			},
		},
	}

	suite.sendSum = baseModels.SendSummary{
		VacancyID:      uint64(7),
		SummaryID:      suite.summary.ID,
		UserID:         suite.person.ID,
		OrganizationID: uint64(13),
		Accepted:       true,
		Denied:         false,
	}

	suite.repo = mock.NewMockSummaryRepository(suite.controller)
	suite.uc = NewSummaryUseCase(suite.repo, bluemonday.UGCPolicy())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestCreateSummary() {
	suite.repo.EXPECT().CreateSummary(&suite.summary).Return(suite.summary.ID, nil).Times(1)

	id, err := suite.uc.CreateSummary(&suite.summary)

	assert.Equal(suite.T(), suite.summary.ID, id)
	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetAllSummaries() {
	suite.repo.EXPECT().GetAllSummaries(3).Return(baseModels.Summaries{}, nil).Times(1)

	_, err := suite.uc.GetAllSummaries("3")

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetUserSummaries() {
	suite.repo.EXPECT().GetUserSummaries(3, uint64(3)).Return(baseModels.Summaries{}, nil).Times(1)

	_, err := suite.uc.GetUserSummaries("3", uint64(3))

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetSummary() {
	suite.repo.EXPECT().GetSummary(uint64(3)).Return(&suite.summary, nil).Times(1)

	_, err := suite.uc.GetSummary(uint64(3))

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestChangeSummary() {
	suite.repo.EXPECT().ChangeSummary(&suite.summary).Return(nil).Times(1)
	suite.repo.EXPECT().CheckAuthor(suite.summary.ID, suite.summary.Author.ID).Return(nil).Times(1)

	err := suite.uc.ChangeSummary(&suite.summary)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestDeleteSummary() {
	suite.repo.EXPECT().DeleteSummary(suite.summary.ID).Return(nil).Times(1)
	suite.repo.EXPECT().CheckAuthor(suite.summary.ID, suite.summary.Author.ID).Return(nil).Times(1)

	err := suite.uc.DeleteSummary(suite.summary.ID, suite.summary.Author.ID)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSendSummary() {
	suite.repo.EXPECT().SendSummary(&suite.sendSum).Return(nil).Times(1)
	suite.repo.EXPECT().CheckAuthor(suite.sendSum.SummaryID, suite.sendSum.UserID).Return(nil).Times(1)

	err := suite.uc.SendSummary(&suite.sendSum)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestSendSummaryRefresh() {
	suite.repo.EXPECT().SendSummary(&suite.sendSum).Return(summaryInterfaces.ErrSummaryAlreadySent).Times(1)
	suite.repo.EXPECT().CheckAuthor(suite.sendSum.SummaryID, suite.sendSum.UserID).Return(nil).Times(1)
	suite.repo.EXPECT().RefreshSummary(suite.sendSum.SummaryID, suite.sendSum.VacancyID).Return(nil).Times(1)

	err := suite.uc.SendSummary(&suite.sendSum)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetOrgSendSummaries() {
	suite.repo.EXPECT().GetOrgSendSummaries(suite.sendSum.OrganizationID).Return(baseModels.OrgSummaries{}, nil).Times(1)

	_, err := suite.uc.GetOrgSendSummaries(suite.sendSum.OrganizationID)

	assert.NoError(suite.T(), err)
}

func (suite *userSuite) TestGetUserSendSummaries() {
	suite.repo.EXPECT().GetUserSendSummaries(suite.sendSum.UserID).Return(baseModels.OrgSummaries{}, nil).Times(1)

	_, err := suite.uc.GetUserSendSummaries(suite.sendSum.UserID)

	assert.NoError(suite.T(), err)
}
