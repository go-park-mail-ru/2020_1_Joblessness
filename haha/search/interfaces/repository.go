package searchInterfaces

import (
	"joblessness/haha/models/base"
)

type SearchRepository interface {
	SearchPersons(params *baseModels.SearchParams) (result []*baseModels.Person, err error)
	SearchOrganizations(params *baseModels.SearchParams) (result []*baseModels.Organization, err error)
	SearchVacancies(params *baseModels.SearchParams) (result []*baseModels.Vacancy, err error)
}
