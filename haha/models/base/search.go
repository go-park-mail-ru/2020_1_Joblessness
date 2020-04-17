package baseModels

import "github.com/microcosm-cc/bluemonday"

type SearchResult struct {
	Persons       []*Person       `json:"persons"`
	Organizations []*Organization `json:"organizations"`
	Vacancies     []*Vacancy      `json:"vacancies"`
}

func (s *SearchResult) Sanitize(policy *bluemonday.Policy) {
	for _, v := range s.Persons {
		v.Sanitize(policy)
	}
	for _, v := range s.Organizations {
		v.Sanitize(policy)
	}
	for _, v := range s.Vacancies {
		v.Sanitize(policy)
	}
}
