package baseModels

import (
	"github.com/microcosm-cc/bluemonday"
	"testing"
)

func TestFlow(t *testing.T) {
	p := &Person{
		ID:        12,
		Login:     "new username",
		Password:  "NewPassword123",
		FirstName: "new first name",
		LastName:  "new last name",
		Email:     "new@email.ru",
		Phone:     "new phone number",
	}
	o := &Organization{
		ID:       12,
		Login:    "new username",
		Password: "NewPassword123",
		Name:     "new name",
		Site:     "new site",
		Email:    "new@email.ru",
		Phone:    "new phone number",
		Tag:      "awdawdawd",
	}
	v := &Vacancy{
		ID: 3,
		Organization: VacancyOrganization{
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

	sr := SearchResult{
		Persons:       Persons{p},
		Organizations: Organizations{o},
		Vacancies:     Vacancies{v},
	}

	sum := &Summary{
		ID: 3,
		Author: Author{
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
		Educations: []Education{
			Education{
				Institution: "was",
				Speciality:  "is",
				Type:        "none",
			},
		},
		Experiences: []Experience{
			Experience{
				CompanyName:      "comp",
				Role:             "role",
				Responsibilities: "response",
			},
		},
	}

	sumResp := &VacancyResponse{
		UserID:    sum.Author.ID,
		Tag:       sum.Author.Tag,
		VacancyID: uint64(7),
		SummaryID: sum.ID,
		Keywords:  sum.Keywords,
	}

	orgSum := OrgSummaries{sumResp}

	sums := Summaries{sum}

	policy := bluemonday.UGCPolicy()

	sr.Sanitize(policy)
	sum.Sanitize(policy)
	sumResp.Sanitize(policy)
	orgSum.Sanitize(policy)
	sums.Sanitize(policy)
}
