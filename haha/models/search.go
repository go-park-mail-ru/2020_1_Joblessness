package models

type SearchResult struct {
	Persons []*Person `json:"persons,omitempty"`
	Organizations []*Organization `json:"organizations,omitempty"`
	Vacancies []*Vacancy `json:"vacancies,omitempty"`
}