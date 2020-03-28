package summaryRepoPostgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"joblessness/haha/models"
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
		PersonID:       1,
		Tag:            "tag",
		Email:          "email",
		Phone:          "phone",
		Registered:     time.Now(),
		Avatar:         "avatar",
	}

	suite.person = Person{
		ID:       uint64(1),
		Name:     "name",
		Gender:   "gender",
		Birthday: time.Now(),
	}
}

func (suite *summarySuite) TearDown() {
	assert.NoError(suite.T(), suite.db.Close())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(summarySuite))
}
