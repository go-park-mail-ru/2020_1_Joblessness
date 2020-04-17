package interviewPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/models/base"
	summaryInterfaces "joblessness/haha/summary/interfaces"
	"testing"
	"time"
)

type summarySuite struct {
	suite.Suite
	rep        *SummaryRepository
	db         *sql.DB
	mock       sqlmock.Sqlmock
	summary    baseModels.Summary
	education  Education
	experience Experience
	user       User
	person     Person
	response   baseModels.VacancyResponse
	sendSum    baseModels.SendSummary
}

func (suite *summarySuite) TestIsOrganizationSummaryTrue() {
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(true)
	suite.mock.
		ExpectQuery("SELECT v.organization_id").
		WithArgs(suite.summary.ID, suite.summary.Author.ID).
		WillReturnRows(rows)

	err := suite.rep.IsOrganizationVacancy(suite.summary.ID, suite.summary.Author.ID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestIsOrganizationSummaryFalse() {
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(false)
	suite.mock.
		ExpectQuery("SELECT v.organization_id").
		WithArgs(suite.summary.ID, suite.summary.Author.ID).
		WillReturnRows(rows)

	err := suite.rep.IsOrganizationVacancy(suite.summary.ID, suite.summary.Author.ID)
	assert.True(suite.T(), errors.Is(err, summaryInterfaces.ErrOrganizationIsNotOwner))
}

func (suite *summarySuite) TestIsOrganizationSummaryFailed() {
	suite.mock.
		ExpectQuery("SELECT v.organization_id").
		WithArgs(suite.summary.ID, suite.summary.Author.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.IsOrganizationVacancy(suite.summary.ID, suite.summary.Author.ID)
	assert.EqualError(suite.T(), err, "")
}

func (suite *summarySuite) TestResponseSummary() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.Accepted, suite.sendSum.Denied, suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.ResponseSummary(&suite.sendSum)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestResponseSummaryNo() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.Accepted, suite.sendSum.Denied, suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := suite.rep.ResponseSummary(&suite.sendSum)
	assert.True(suite.T(), errors.Is(err, summaryInterfaces.ErrNoSummaryToRefresh))
}

func (suite *summarySuite) TestResponseSummaryFailed() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnError(errors.New(""))

	err := suite.rep.RefreshSummary(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.Error(suite.T(), err)
}