package summaryRepoPostgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/models"
	summaryInterfaces "joblessness/haha/summary/interfaces"
	"testing"
	"time"
)

type summarySuite struct {
	suite.Suite
	rep *SummaryRepository
	db *sql.DB
	mock sqlmock.Sqlmock
	summary models.Summary
	education Education
	experience Experience
	user User
	person Person
	response models.VacancyResponse
	sendSum models.SendSummary
}

func (suite *summarySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.rep = NewSummaryRepository(suite.db)

	suite.summary = models.Summary{
		ID:          3,
		Author:      models.Author{
			ID:        12,
			Tag:       "tag",
			Email:     "email",
			Phone:     "phone",
			Avatar:    "avatar",
			FirstName: "first",
			LastName:  "name",
			Gender:    "gender",
		},
		Keywords:    "key",
		Educations:  []models.Education{
			models.Education{
				Institution: "was",
				Speciality:  "is",
				Type:        "none",
			},
		},
		Experiences: []models.Experience{
			models.Experience{
				CompanyName:      "comp",
				Role:             "role",
				Responsibilities: "response",
				Start:            time.Now(),
				Stop:             time.Now().AddDate(1, 1, 1),
			},
		},
	}

	suite.user = User{
		ID:             12,
		OrganizationID: 0,
		PersonID:       3,
		Tag:            "tag",
		Email:          "email",
		Phone:          "phone",
		Registered:     time.Now(),
		Avatar:         "avatar",
	}

	suite.person = Person{
		ID:       uint64(3),
		Name:     "name",
		Gender:   "gender",
		Birthday: time.Now(),
	}

	suite.response = models.VacancyResponse{
		UserID:    suite.person.ID,
		Tag:       suite.user.Tag,
		VacancyID: uint64(7),
		SummaryID: suite.summary.ID,
		Keywords:  suite.summary.Keywords,
	}

	suite.sendSum = models.SendSummary{
		VacancyID:      uint64(7),
		SummaryID:      suite.summary.ID,
		UserID:    		suite.person.ID,
		OrganizationID: uint64(13),
		Accepted:       true,
		Denied:         false,
	}
}

func (suite *summarySuite) TearDown() {
	assert.NoError(suite.T(), suite.db.Close())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(summarySuite))
}

func (suite *summarySuite) TestCreateSummary() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(uint64(3))

	suite.mock.
		ExpectQuery("INSERT INTO summary").
		WithArgs(suite.summary.Author.ID, suite.summary.Keywords, suite.summary.Name, suite.summary.SalaryFrom,
			suite.summary.SalaryTo).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO education").
		WithArgs(suite.summary.ID, suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
		nil , suite.summary.Educations[0].Type).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec("INSERT INTO experience").
		WithArgs(suite.summary.ID, suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
		suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start, suite.summary.Experiences[0].Stop).
		WillReturnResult(sqlmock.NewResult(1, 1))

	summaryID, err := suite.rep.CreateSummary(&suite.summary)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.summary.ID, summaryID)
}

func (suite *summarySuite) TestCreateSummaryFailedOne() {
	suite.mock.
		ExpectQuery("INSERT INTO summary").
		WithArgs(suite.summary.Author.ID, suite.summary.Keywords).
		WillReturnError(errors.New(""))

	_, err := suite.rep.CreateSummary(&suite.summary)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestCreateSummaryFailedTwo() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(uint64(3))

	suite.mock.
		ExpectQuery("INSERT INTO summary").
		WithArgs(suite.summary.Author.ID, suite.summary.Keywords).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO education").
		WithArgs(suite.summary.ID, suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			nil , suite.summary.Educations[0].Type).
		WillReturnError(errors.New(""))

	_, err := suite.rep.CreateSummary(&suite.summary)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestCreateSummaryFailedThree() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(uint64(3))

	suite.mock.
		ExpectQuery("INSERT INTO summary").
		WithArgs(suite.summary.Author.ID, suite.summary.Keywords).
		WillReturnRows(rows)
	suite.mock.
		ExpectExec("INSERT INTO education").
		WithArgs(suite.summary.ID, suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			nil , suite.summary.Educations[0].Type).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec("INSERT INTO experience").
		WithArgs(suite.summary.ID, suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
			suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start, suite.summary.Experiences[0].Stop).
		WillReturnError(errors.New(""))

	_, err := suite.rep.CreateSummary(&suite.summary)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetSummary() {
	rows := sqlmock.NewRows([]string{"author", "keywords"}).
		AddRow(suite.summary.Author.ID, suite.summary.Keywords)
	suite.mock.
		ExpectQuery("SELECT author, keywords").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"institution", "speciality", "graduated", "type"}).
		AddRow(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			suite.summary.Educations[0].Graduated, suite.summary.Educations[0].Type)
	suite.mock.
		ExpectQuery("SELECT institution, speciality, graduated, type").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"company_name", "role", "responsibilities", "start", "stop"}).
		AddRow(suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
			suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start,
			suite.summary.Experiences[0].Stop)
	suite.mock.
		ExpectQuery("SELECT company_name, role, responsibilities, start, stop").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"tag", "email", "phone", "avatar", "name", "gender", "birthday"}).
		AddRow(suite.summary.Author.Tag, suite.summary.Author.Email, suite.summary.Author.Phone,
			suite.summary.Author.Avatar, "first name", suite.summary.Author.Gender, suite.summary.Author.Birthday)
	suite.mock.
		ExpectQuery("SELECT tag, email, phone, avatar, name, gender, birthday").
		WithArgs(suite.summary.Author.ID).
		WillReturnRows(rows)

	summary, err := suite.rep.GetSummary(suite.summary.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.summary, *summary)
}

func (suite *summarySuite) TestGetSummaryFailedOne() {
	suite.mock.
		ExpectQuery("SELECT author, keywords").
		WithArgs(suite.summary.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetSummary(suite.summary.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetSummaryFailedTwo() {
	rows := sqlmock.NewRows([]string{"author", "keywords"}).
		AddRow(suite.summary.Author.ID, suite.summary.Keywords)
	suite.mock.
		ExpectQuery("SELECT author, keywords").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT institution, speciality, graduated, type").
		WithArgs(suite.summary.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetSummary(suite.summary.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetSummaryFailedThree() {
	rows := sqlmock.NewRows([]string{"author", "keywords"}).
		AddRow(suite.summary.Author.ID, suite.summary.Keywords)
	suite.mock.
		ExpectQuery("SELECT author, keywords").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"institution", "speciality", "graduated", "type"}).
		AddRow(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			suite.summary.Educations[0].Graduated, suite.summary.Educations[0].Type)
	suite.mock.
		ExpectQuery("SELECT institution, speciality, graduated, type").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT company_name, role, responsibilities, start, stop").
		WithArgs(suite.summary.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetSummary(suite.summary.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetSummaryFailedFor() {
	rows := sqlmock.NewRows([]string{"author", "keywords"}).
		AddRow(suite.summary.Author.ID, suite.summary.Keywords)
	suite.mock.
		ExpectQuery("SELECT author, keywords").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"institution", "speciality", "graduated", "type"}).
		AddRow(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			suite.summary.Educations[0].Graduated, suite.summary.Educations[0].Type)
	suite.mock.
		ExpectQuery("SELECT institution, speciality, graduated, type").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"company_name", "role", "responsibilities", "start", "stop"}).
		AddRow(suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
			suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start,
			suite.summary.Experiences[0].Stop)
	suite.mock.
		ExpectQuery("SELECT company_name, role, responsibilities, start, stop").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	suite.mock.
		ExpectQuery("SELECT tag, email, phone, avatar, name, gender, birthday").
		WithArgs(suite.summary.Author.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetSummary(suite.summary.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetSummaries() {
	rows := sqlmock.NewRows([]string{"id", "author", "keywords", "name", "salary_from", "salary_to"}).
		AddRow(suite.summary.ID, suite.summary.Author.ID, suite.summary.Keywords, suite.summary.Name,
			suite.summary.SalaryFrom, suite.summary.SalaryTo)
	suite.mock.
		ExpectQuery("SELECT id, author, keywords").
		WithArgs(10, uint64(10)).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"institution", "speciality", "graduated", "type"}).
		AddRow(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			suite.summary.Educations[0].Graduated, suite.summary.Educations[0].Type)
	suite.mock.
		ExpectQuery("SELECT institution, speciality, graduated, type").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"company_name", "role", "responsibilities", "start", "stop"}).
		AddRow(suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
			suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start,
			suite.summary.Experiences[0].Stop)
	suite.mock.
		ExpectQuery("SELECT company_name, role, responsibilities, start, stop").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"tag", "email", "phone", "avatar", "name", "gender", "birthday"}).
		AddRow(suite.summary.Author.Tag, suite.summary.Author.Email, suite.summary.Author.Phone,
			suite.summary.Author.Avatar, "first name", suite.summary.Author.Gender, suite.summary.Author.Birthday)
	suite.mock.
		ExpectQuery("SELECT tag, email, phone, avatar, name, gender, birthday").
		WithArgs(suite.summary.Author.ID).
		WillReturnRows(rows)

	summaries, err := suite.rep.GetAllSummaries(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.summary, summaries[0])
}

func (suite *summarySuite) TestGetUserSummaries() {
	rows := sqlmock.NewRows([]string{"id", "author", "keywords", "name", "salary_from", "salary_to"}).
		AddRow(suite.summary.ID, suite.summary.Author.ID, suite.summary.Keywords, suite.summary.Name,
			suite.summary.SalaryFrom, suite.summary.SalaryTo)
	suite.mock.
		ExpectQuery("SELECT id, author, keywords").
		WithArgs(suite.summary.Author.ID, 10, uint64(0)).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"institution", "speciality", "graduated", "type"}).
		AddRow(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			suite.summary.Educations[0].Graduated, suite.summary.Educations[0].Type)
	suite.mock.
		ExpectQuery("SELECT institution, speciality, graduated, type").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"company_name", "role", "responsibilities", "start", "stop"}).
		AddRow(suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
			suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start,
			suite.summary.Experiences[0].Stop)
	suite.mock.
		ExpectQuery("SELECT company_name, role, responsibilities, start, stop").
		WithArgs(suite.summary.ID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"tag", "email", "phone", "avatar", "name", "gender", "birthday"}).
		AddRow(suite.summary.Author.Tag, suite.summary.Author.Email, suite.summary.Author.Phone,
			suite.summary.Author.Avatar, "first name", suite.summary.Author.Gender, suite.summary.Author.Birthday)
	suite.mock.
		ExpectQuery("SELECT tag, email, phone, avatar, name, gender, birthday").
		WithArgs(suite.summary.Author.ID).
		WillReturnRows(rows)

	summaries, err := suite.rep.GetUserSummaries(0, uint64(12))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.summary, summaries[0])
}

func (suite *summarySuite) TestChangeSummary() {
	suite.mock.
		ExpectExec("UPDATE summary").
		WithArgs(suite.summary.Keywords, suite.summary.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.
		ExpectExec("UPDATE education").
		WithArgs(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
		nil, suite.summary.Educations[0].Type, suite.summary.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.
		ExpectExec("UPDATE experience").
		WithArgs(suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
			suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start,
			suite.summary.Experiences[0].Stop, suite.summary.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.ChangeSummary(&suite.summary)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestChangeSummaryFailedOne() {
	suite.mock.
		ExpectExec("UPDATE summary").
		WithArgs(suite.summary.Keywords, suite.summary.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangeSummary(&suite.summary)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestChangeSummaryFailedTwo() {
	suite.mock.
		ExpectExec("UPDATE summary").
		WithArgs(suite.summary.Keywords, suite.summary.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.
		ExpectExec("UPDATE education").
		WithArgs(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			nil, suite.summary.Educations[0].Type, suite.summary.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangeSummary(&suite.summary)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestChangeSummaryFailedThree() {
	suite.mock.
		ExpectExec("UPDATE summary").
		WithArgs(suite.summary.Keywords, suite.summary.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.
		ExpectExec("UPDATE education").
		WithArgs(suite.summary.Educations[0].Institution, suite.summary.Educations[0].Speciality,
			nil, suite.summary.Educations[0].Type, suite.summary.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.
		ExpectExec("UPDATE experience").
		WithArgs(suite.summary.Experiences[0].CompanyName, suite.summary.Experiences[0].Role,
			suite.summary.Experiences[0].Responsibilities, suite.summary.Experiences[0].Start,
			suite.summary.Experiences[0].Stop, suite.summary.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.ChangeSummary(&suite.summary)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestDeleteSummary() {
	suite.mock.
		ExpectExec("DELETE FROM summary").
		WithArgs(suite.summary.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.DeleteSummary(suite.summary.ID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestDeleteSummaryFailed() {
	suite.mock.
		ExpectExec("DELETE FROM summary").
		WithArgs(suite.summary.ID).
		WillReturnError(errors.New(""))

	err := suite.rep.DeleteSummary(suite.summary.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestIsPersonSummaryTrue() {
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(suite.summary.Author.ID)
	suite.mock.
		ExpectQuery("SELECT u.id").
		WithArgs(suite.summary.ID, suite.summary.Author.ID).
		WillReturnRows(rows)

	res, err := suite.rep.IsPersonSummary(suite.summary.ID, suite.summary.Author.ID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), res)
}

func (suite *summarySuite) TestIsPersonSummaryFalse() {
	rows := sqlmock.NewRows([]string{"id"})
	suite.mock.
		ExpectQuery("SELECT u.id").
		WithArgs(suite.summary.ID, suite.summary.Author.ID).
		WillReturnRows(rows)

	res, err := suite.rep.IsPersonSummary(suite.summary.ID, suite.summary.Author.ID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), res)
}

func (suite *summarySuite) TestIsPersonSummaryFailed() {
	suite.mock.
		ExpectQuery("SELECT u.id").
		WithArgs(suite.summary.Author.ID, suite.summary.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.IsPersonSummary(suite.summary.ID, suite.summary.Author.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestIsOrganizationSummaryTrue() {
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(suite.summary.Author.ID)
	suite.mock.
		ExpectQuery("SELECT u.id").
		WithArgs(suite.summary.ID, suite.summary.Author.ID).
		WillReturnRows(rows)

	res, err := suite.rep.IsOrganizationVacancy(suite.summary.ID, suite.summary.Author.ID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), res)
}

func (suite *summarySuite) TestIsOrganizationSummaryFalse() {
	rows := sqlmock.NewRows([]string{"id"})
	suite.mock.
		ExpectQuery("SELECT u.id").
		WithArgs(suite.summary.ID, suite.summary.Author.ID).
		WillReturnRows(rows)

	res, err := suite.rep.IsOrganizationVacancy(suite.summary.ID, suite.summary.Author.ID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), res)
}

func (suite *summarySuite) TestIsOrganizationSummaryFailed() {
	suite.mock.
		ExpectQuery("SELECT u.id").
		WithArgs(suite.summary.Author.ID, suite.summary.ID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.IsOrganizationVacancy(suite.summary.ID, suite.summary.Author.ID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestSendSummary() {
	suite.mock.
		ExpectExec("INSERT INTO response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.SendSummary(&suite.sendSum)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestSendSummaryAlreadySend() {
	suite.mock.
		ExpectExec("INSERT INTO response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := suite.rep.SendSummary(&suite.sendSum)
	assert.Equal(suite.T(), summaryInterfaces.ErrSummaryAlreadySend, err)
}

func (suite *summarySuite) TestSendSummaryFailed() {
	suite.mock.
		ExpectExec("INSERT INTO response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnError(errors.New(""))

	err := suite.rep.SendSummary(&suite.sendSum)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestRefreshSummary() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.rep.RefreshSummary(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestRefreshSummaryNoSummary() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := suite.rep.RefreshSummary(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.Equal(suite.T(), summaryInterfaces.ErrNoSummaryToRefresh, err)
}

func (suite *summarySuite) TestRefreshSummaryFailed() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnError(errors.New(""))

	err := suite.rep.RefreshSummary(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.Error(suite.T(), err)
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
	assert.Equal(suite.T(), summaryInterfaces.ErrNoSummaryToRefresh, err)
}

func (suite *summarySuite) TestResponseSummaryFailed() {
	suite.mock.
		ExpectExec("UPDATE response").
		WithArgs(suite.sendSum.SummaryID, suite.sendSum.VacancyID).
		WillReturnError(errors.New(""))

	err := suite.rep.RefreshSummary(suite.sendSum.SummaryID, suite.sendSum.VacancyID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetOrgSummaries() {
	rows := sqlmock.NewRows([]string{"id", "tag", "id", "id", "keywords", "name", "name", "approved", "rejected"}).
		AddRow(suite.response.UserID, suite.response.Tag, suite.response.VacancyID, suite.response.SummaryID,
			suite.response.Keywords, suite.response.SummaryName, suite.response.VacancyName, suite.response.Accepted,
			suite.response.Denied)

	suite.mock.
		ExpectQuery("SELECT u.id, u.tag, v.id, s.id, s.keywords").
		WithArgs(suite.sendSum.OrganizationID).
		WillReturnRows(rows)

	res, err := suite.rep.GetOrgSendSummaries(suite.sendSum.OrganizationID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.response, *res[0])
}

func (suite *summarySuite) TestGetOrgSummariesFailed() {
	suite.mock.
		ExpectQuery("SELECT u.id, u.tag, v.id, s.id, s.keywords").
		WithArgs(suite.sendSum.OrganizationID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetOrgSendSummaries(suite.sendSum.OrganizationID)
	assert.Error(suite.T(), err)
}

func (suite *summarySuite) TestGetUserSendSummaries() {
	rows := sqlmock.NewRows([]string{"id", "id", "keywords", "name", "name", "approved", "rejected"}).
		AddRow(suite.response.VacancyID, suite.response.SummaryID,
			suite.response.Keywords, suite.response.SummaryName, suite.response.VacancyName, suite.response.Accepted,
			suite.response.Denied)

	suite.mock.
		ExpectQuery("SELECT v.id, s.id, s.keywords, s.name, v.name").
		WithArgs(suite.sendSum.OrganizationID).
		WillReturnRows(rows)

	_, err := suite.rep.GetUserSendSummaries(suite.sendSum.OrganizationID)
	assert.NoError(suite.T(), err)
}

func (suite *summarySuite) TestGetUserSendSummariesFailed() {
	suite.mock.
		ExpectQuery("SELECT v.id, s.id, s.keywords, s.name, v.name").
		WithArgs(suite.sendSum.OrganizationID).
		WillReturnError(errors.New(""))

	_, err := suite.rep.GetUserSendSummaries(suite.sendSum.OrganizationID)
	assert.Error(suite.T(), err)
}