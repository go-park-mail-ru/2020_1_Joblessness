package models

type SearchResult struct {
	Persons []*Person `json:"persons"`
	Organizations []*Organization `json:"organizations"`
	Vacancies []*Vacancy `json:"vacancies"`
}