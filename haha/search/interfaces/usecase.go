package searchInterfaces

import (
	"joblessness/haha/models/base"
)

type SearchUseCase interface {
	Search(searchType, request, since, desc string) (result baseModels.SearchResult, err error)
}
