package searchUseCase

import (
	"github.com/microcosm-cc/bluemonday"
	"joblessness/haha/models/base"
	searchInterfaces "joblessness/haha/search/interfaces"
)

type SearchUseCase struct {
	searchRepo searchInterfaces.SearchRepository
	policy     *bluemonday.Policy
}

func NewSearchUseCase(userRepo searchInterfaces.SearchRepository, policy *bluemonday.Policy) *SearchUseCase {
	return &SearchUseCase{
		searchRepo: userRepo,
		policy:     policy,
	}
}

func (a *SearchUseCase) Search(searchType, request, since, desc string) (result baseModels.SearchResult, err error) {
	params := &baseModels.SearchParams{
		Request: request,
		Since:   since,
		Desc:    desc,
	}

	switch searchType {
	case "person":
		result.Persons, err = a.searchRepo.SearchPersons(params)
	case "organization":
		result.Organizations, err = a.searchRepo.SearchOrganizations(params)
	case "vacancy":
		result.Vacancies, err = a.searchRepo.SearchVacancies(params)
	case "":
		result.Persons, err = a.searchRepo.SearchPersons(params)
		if err != nil {
			break
		}
		result.Organizations, err = a.searchRepo.SearchOrganizations(params)
		if err != nil {
			break
		}
		result.Vacancies, err = a.searchRepo.SearchVacancies(params)
	default:
		return result, searchInterfaces.ErrUnknownRequest
	}

	result.Sanitize(a.policy)
	return result, err
}
