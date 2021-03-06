package interviewPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	"joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
	"log"
	"testing"
	"time"
)

type summarySuite struct {
	suite.Suite
	rep       *InterviewRepository
	db        *sql.DB
	mock      sqlmock.Sqlmock
	sendSum   baseModels.SendSummary
	message   chat.Message
	convTitle baseModels.ConversationTitle
	params    baseModels.ChatParameters
}

func (suite *summarySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewInterviewRepository(suite.db)

	suite.sendSum = baseModels.SendSummary{
		VacancyID:      uint64(7),
		SummaryID:      uint64(2),
		UserID:         uint64(1),
		OrganizationID: uint64(13),
		InterviewDate:  time.Now(),
		Accepted:       true,
		Denied:         false,
	}

	suite.message = chat.Message{
		Message:   "message",
		UserOneID: 1,
		UserOne:   "awd",
		UserTwoID: 2,
		UserTwo:   "awd",
		Created:   time.Now(),
	}

	suite.convTitle = baseModels.ConversationTitle{
		ChatterID:     uint64(1),
		Avatar:        "awd",
		ChatterName:   "awd",
		Tag:           "ad",
		InterviewDate: time.Now(),
	}

	suite.params = baseModels.ChatParameters{
		From: 1,
		To:   2,
		Page: 0,
	}
}

func (suite *summarySuite) TearDown() {
	assert.NoError(suite.T(), suite.db.Close())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(summarySuite))
}

func (suite *summarySuite) TestIsOrganizationSummaryTrue() {
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(true)
	suite.mock.
		ExpectQuery("SELECT v.organization_id").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.OrganizationID).
		WillReturnRows(rows)

	err := suite.rep.IsOrganizationVacancy(suite.sendSum.SummaryID, suite.sendSum.OrganizationID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestIsOrganizationSummaryFalse() {
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(false)
	suite.mock.
		ExpectQuery("SELECT v.organization_id").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.OrganizationID).
		WillReturnRows(rows)

	err := suite.rep.IsOrganizationVacancy(suite.sendSum.SummaryID, suite.sendSum.OrganizationID)
	assert.True(suite.T(), errors.Is(err, interviewInterfaces.ErrOrganizationIsNotOwner))
}

func (suite *summarySuite) TestIsOrganizationSummaryFailed() {
	suite.mock.
		ExpectQuery("SELECT v.organization_id").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.OrganizationID).
		WillReturnError(errors.New(""))

	err := suite.rep.IsOrganizationVacancy(suite.sendSum.SummaryID, suite.sendSum.OrganizationID)
	assert.EqualError(suite.T(), err, "")
}

func (suite *summarySuite) TestResponseSummary() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.Accepted, suite.sendSum.Denied, suite.sendSum.InterviewDate,
			suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.ResponseSummary(&suite.sendSum)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestResponseSummaryNo() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.Accepted, suite.sendSum.Denied, suite.sendSum.InterviewDate,
			suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := suite.rep.ResponseSummary(&suite.sendSum)
	log.Println("ERROR: ", err)
	assert.True(suite.T(), errors.Is(err, interviewInterfaces.ErrNoSummaryToRefresh))
}

func (suite *summarySuite) TestSaveMessage() {
	suite.mock.
		ExpectExec("INSERT INTO message").
		WithArgs(suite.message.UserOneID, suite.message.UserTwoID, suite.message.UserOne, suite.message.UserTwo,
			suite.message.Message).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := suite.rep.SaveMessage(&suite.message)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestSaveMessageFailed() {
	suite.mock.
		ExpectExec("INSERT INTO message").
		WithArgs(suite.message.UserOneID, suite.message.UserTwoID, suite.message.UserOne, suite.message.UserTwo,
			suite.message.Message).
		WillReturnError(errors.New(""))

	err := suite.rep.SaveMessage(&suite.message)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetHistory() {
	rows := sqlmock.NewRows([]string{"user_one_id", "user_two_id", "user_one", "user_two", "body", "created"}).
		AddRow(suite.message.UserOneID, suite.message.UserTwoID, suite.message.UserOne, suite.message.UserTwo,
			suite.message.Message, suite.message.Created)
	suite.mock.
		ExpectQuery("SELECT user_one_id, user_two_id").
		WithArgs(suite.message.UserOneID, suite.message.UserTwoID, uint64(40), uint64(0)).
		WillReturnRows(rows)

	_, err := suite.rep.GetHistory(&suite.params)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestGetHistoryFailed() {
	suite.mock.
		ExpectQuery("SELECT user_one_id, user_two_id").
		WithArgs(suite.message.UserOneID, suite.message.UserTwoID, uint64(20), uint64(0)).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetHistory(&suite.params)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetUserSendMessages() {
	rows := sqlmock.NewRows([]string{"user_one_id", "user_two_id", "user_one", "user_two", "body", "created"}).
		AddRow(suite.message.UserOneID, suite.message.UserTwoID, suite.message.UserOne, suite.message.UserTwo,
			suite.message.Message, suite.message.Created)
	suite.mock.
		ExpectQuery("SELECT user_one_id, user_two_id").
		WithArgs(suite.message.UserOneID, suite.message.UserTwoID, uint64(40), uint64(0)).
		WillReturnRows(rows)

	_, _, err := suite.rep.getUserSendMessages(&suite.params)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestGetUserSendMessagesFailed() {
	suite.mock.
		ExpectQuery("SELECT user_one_id, user_two_id").
		WithArgs(suite.message.UserOneID, suite.message.UserTwoID, uint64(20), uint64(0)).
		WillReturnError(errors.New(""))

	_, _, err := suite.rep.getUserSendMessages(&suite.params)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetResponseCredentials() {
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(suite.message.UserTwoID, suite.message.UserTwo)
	suite.mock.
		ExpectQuery("SELECT u.id, p.name").
		WithArgs(suite.sendSum.SummaryID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "name"}).
		AddRow(suite.message.UserOneID, suite.message.UserOne)
	suite.mock.
		ExpectQuery("SELECT u.id, o.name").
		WithArgs(suite.sendSum.VacancyID).
		WillReturnRows(rows)

	_, err := suite.rep.GetResponseCredentials(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestGetResponseCredentialsFailedOne() {
	suite.mock.
		ExpectQuery("SELECT u.id, p.name").
		WithArgs(suite.sendSum.SummaryID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetResponseCredentials(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetResponseCredentialsFailedTwo() {
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(suite.message.UserTwoID, suite.message.UserTwo)
	suite.mock.
		ExpectQuery("SELECT u.id, p.name").
		WithArgs(suite.sendSum.SummaryID).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT u.id, o.name").
		WithArgs(suite.sendSum.VacancyID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetResponseCredentials(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetConversations() {
	rows := sqlmock.NewRows([]string{"id", "avatar", "name", "tag"}).
		AddRow(suite.convTitle.ChatterID, suite.convTitle.Avatar, suite.convTitle.ChatterName,
			suite.convTitle.InterviewDate)
	suite.mock.
		ExpectQuery("SELECT u_to.id, u_to.avatar").
		WithArgs(suite.message.UserOneID).
		WillReturnRows(rows)

	_, err := suite.rep.GetConversations(suite.message.UserOneID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestGetConversationsFailed() {
	suite.mock.
		ExpectQuery("SELECT u.id, u.tag, r.interview_date").
		WithArgs(suite.message.UserOneID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetConversations(suite.message.UserOneID)
	assert.Error(suite.T(), err)
}
