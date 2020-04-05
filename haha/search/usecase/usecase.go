package searchUseCase

import (
	"joblessness/haha/models"
	searchInterfaces "joblessness/haha/search/interfaces"
)

type SearchUseCase struct {
	searchRepo searchInterfaces.SearchRepository
}

func NewSearchUseCase(userRepo searchInterfaces.SearchRepository) *SearchUseCase {
	return &SearchUseCase{
		searchRepo:userRepo,
	}
}

func (a *SearchUseCase) Search(searchType, request, since, desc string) (result models.SearchResult, err error) {
	switch searchType {
	case "person":
		result.Persons, err = a.searchRepo.SearchPersons(request, since, desc)
	case "organization":
		result.Organizations, err = a.searchRepo.SearchOrganizations(request, since, desc)
	case "vacancy":
		result.Vacancies, err = a.searchRepo.SearchVacancies(request, since, desc)
	case "":
		result.Persons, err = a.searchRepo.SearchPersons(request, since, desc)
		if err != nil {
			break
		}
		result.Organizations, err = a.searchRepo.SearchOrganizations(request, since, desc)
		if err != nil {
			break
		}
		result.Vacancies, err = a.searchRepo.SearchVacancies(request, since, desc)
	default:
		return result, searchInterfaces.ErrUnknownRequest
	}

	return result, err
}