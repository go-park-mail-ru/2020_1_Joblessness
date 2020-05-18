package grpcModels

import (
	"github.com/microcosm-cc/bluemonday"
	baseModels "joblessness/haha/models/base"
	"testing"
)

func TestFlow(t *testing.T) {
	p := &baseModels.Person{
		ID:        12,
		Login:     "new username",
		Password:  "NewPassword123",
		FirstName: "new first name",
		LastName:  "new last name",
		Email:     "new@email.ru",
		Phone:     "new phone number",
	}
	o := &baseModels.Organization{
		ID:       12,
		Login:    "new username",
		Password: "NewPassword123",
		Name:     "new name",
		Site:     "new site",
		Email:    "new@email.ru",
		Phone:    "new phone number",
		Tag:      "awdawdawd",
	}
	v := &baseModels.Vacancy{
		ID: 3,
		Organization: baseModels.VacancyOrganization{
			ID:     12,
			Tag:    "",
			Email:  "",
			Phone:  "",
			Avatar: "",
			Name:   "",
			Site:   "",
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

	sr := baseModels.SearchResult{
		Persons:       baseModels.Persons{p},
		Organizations: baseModels.Organizations{o},
		Vacancies:     baseModels.Vacancies{v},
	}

	sum := &baseModels.Summary{
		ID: 3,
		Author: baseModels.Author{
			ID:        12,
			Tag:       "tag",
			Email:     "email",
			Phone:     "phone",
			Avatar:    "avatar",
			FirstName: "first",
			LastName:  "name",
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
			},
		},
	}

	sumResp := &baseModels.VacancyResponse{
		UserID:    sum.Author.ID,
		Tag:       sum.Author.Tag,
		VacancyID: uint64(7),
		SummaryID: sum.ID,
		Avatar:  "avatar",
	}

	orgSum := baseModels.OrgSummaries{sumResp}

	sums := baseModels.Summaries{sum}

	policy := bluemonday.UGCPolicy()

	sr.Sanitize(policy)
	sum.Sanitize(policy)
	sumResp.Sanitize(policy)
	orgSum.Sanitize(policy)
	sums.Sanitize(policy)
}
