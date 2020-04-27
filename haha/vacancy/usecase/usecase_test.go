package vacancyUseCase

import (
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/suite"
	baseModels "joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
	mockRoom "joblessness/haha/utils/chat/mock"
	"joblessness/haha/vacancy/repository/mock"
	"testing"
	"time"
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
	response   baseModels.VacancyResponse
	sendSum    baseModels.SendSummary
	room *mockRoom.MockRoom
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

	suite.room = mockRoom.NewMockRoom(suite.controller)

	suite.repo = mock.NewMockVacancyRepository(suite.controller)
	suite.uc = NewVacancyUseCase(suite.repo, suite.room, bluemonday.UGCPolicy())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestCreateSummary() {
	//suite.repo.EXPECT().CreateSummary(&suite.summary).Return(suite.summary.ID, nil).Times(1)
	//
	//id, err := suite.uc.CreateSummary(&suite.summary)
	//
	//assert.Equal(suite.T(), suite.summary.ID, id)
	//assert.NoError(suite.T(), err)
}