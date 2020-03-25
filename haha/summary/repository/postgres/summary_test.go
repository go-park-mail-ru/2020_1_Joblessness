package summaryRepoPostgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/models"
	"testing"
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
			Birthday:  "birth",
		},
		Keywords:    "key",
		Educations:  []models.Education{
			models.Education{
				Institution: "was",
				Speciality:  "is",
				Graduated:   "yes",
				Type:        "none",
			},
		},
		Experiences: []models.Experience{
			models.Experience{
				CompanyName:      "comp",
				Role:             "role",
				Responsibilities: "response",
				Start:            "start",
				Stop:             "stop",
			},
		},
	}

	suite.user = User{
		ID:             12,
		OrganizationID: 0,
		PersonID:       5,
		Tag:            sql.NullString{},
		Email:          sql.NullString{},
		Phone:          sql.NullString{},
		Registered:     sql.NullTime{},
		Avatar:         sql.NullString{},
	}

	suite.person = Person{
		ID:       sql.NullString{},
		Name:     sql.NullString{},
		Gender:   sql.NullString{},
		Birthday: sql.NullTime{},
	}
}

func (suite *summarySuite) TearDown() {
	assert.NoError(suite.T(), suite.db.Close())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(summarySuite))
}
