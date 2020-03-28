package models

type SearchResult struct {
	Persons []*Person
	Organizations []*Organization
	Vacancies []*Vacancy
}