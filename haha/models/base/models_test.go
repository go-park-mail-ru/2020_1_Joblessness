package baseModels

import (
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	res, err := easyjson.Marshal(p)
	assert.NoError(t, err)
	err = easyjson.Unmarshal(res, p)
	assert.NoError(t, err)
}

func TestUserLogin(t *testing.T) {
	userLogin := UserLogin{
		Login:    "login",
		Password: "user",
	}

	res, err := easyjson.Marshal(userLogin)
	assert.NoError(t, err)
	err = easyjson.Unmarshal(res, &userLogin)
	assert.NoError(t, err)
}

func TestResponseRole(t *testing.T) {
	responseRole := ResponseRole{
		ID:   1,
		Role: "awdaw",
	}

	res, err := easyjson.Marshal(responseRole)
	assert.NoError(t, err)
	err = easyjson.Unmarshal(res, &responseRole)
	assert.NoError(t, err)
}

func TestPersons(t *testing.T) {
	persons := Persons{
		&Person{
			ID:         0,
			Login:      "",
			Password:   "",
			Tag:        "",
			Email:      "",
			Phone:      "",
			Registered: time.Now(),
			Avatar:     "",
			FirstName:  "",
			LastName:   "",
			Gender:     "",
			Birthday:   time.Now(),
		},
		&Person{
			ID:         1,
			Login:      "Loawd",
			Password:   "awda",
			Tag:        "adaw",
			Email:      "adwa",
			Phone:      "awd",
			Registered: time.Time{},
			Avatar:     "adaw",
			FirstName:  "awdaw",
			LastName:   "adaw",
			Gender:     "awd",
			Birthday:   time.Time{},
		},
	}

	res, err := easyjson.Marshal(persons)
	assert.NoError(t, err)
	err = easyjson.Unmarshal(res, &persons)
	assert.NoError(t, err)
}

func TestOrgs(t *testing.T) {
	orgs := Organizations{
		&Organization{
			ID:         1,
			Login:      "adaw",
			Password:   "awda",
			Tag:        "adaw",
			Email:      "awda",
			Phone:      "awdaw",
			Registered: time.Time{},
			Avatar:     "awda",
			Name:       "adaw",
			About:      "dawda",
			Site:       "dawda",
		},
		&Organization{
			ID:         0,
			Login:      "",
			Password:   "",
			Tag:        "",
			Email:      "",
			Phone:      "",
			Registered: time.Time{},
			Avatar:     "",
			Name:       "",
			About:      "",
			Site:       "",
		},
	}

	res, err := easyjson.Marshal(orgs)
	assert.NoError(t, err)
	err = easyjson.Unmarshal(res, &orgs)
	assert.NoError(t, err)
}

func TestFavs(t *testing.T) {
	fav := Favorites{
		&Favorite{
			ID:       0,
			Tag:      "",
			Avatar:   "",
			IsPerson: false,
			Name:     "",
			Surname:  "",
		},
		&Favorite{
			ID:       1,
			Tag:      "adaw",
			Avatar:   "daww",
			IsPerson: true,
			Name:     "awda",
			Surname:  "awdaw",
		},
	}

	res, err := easyjson.Marshal(fav)
	assert.NoError(t, err)
	err = easyjson.Unmarshal(res, &fav)
	assert.NoError(t, err)
}