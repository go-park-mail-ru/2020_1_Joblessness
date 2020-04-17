package searchInterfaces

import (
	"joblessness/haha/models/base"
)

type SearchRepository interface {
	SearchPersons(request, since, desc string) (result []*baseModels.Person, err error)
	SearchOrganizations(request, since, desc string) (result []*baseModels.Organization, err error)
	SearchVacancies(request, since, desc string) (result []*baseModels.Vacancy, err error)
}
