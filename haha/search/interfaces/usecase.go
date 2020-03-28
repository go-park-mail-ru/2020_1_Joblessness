package searchInterfaces

import "joblessness/haha/models"

type SearchUseCase interface {
	Search(searchType, request, since, desc string) (result models.SearchResult, err error)
}