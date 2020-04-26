package baseModels

import "github.com/microcosm-cc/bluemonday"

//easyjson:json
type SearchResult struct {
	Persons       Persons      `json:"persons"`
	Organizations Organizations `json:"organizations"`
	Vacancies     Vacancies      `json:"vacancies"`
}

func (s *SearchResult) Sanitize(policy *bluemonday.Policy) {
	s.Persons.Sanitize(policy)
	s.Organizations.Sanitize(policy)
	s.Vacancies.Sanitize(policy)
}

type SearchParams struct {
	Request string
	Since   string
	Desc    string
}
