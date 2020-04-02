package searchInterfaces

import "joblessness/haha/models"

type SearchRepository interface {
	SearchPersons(request, since, desc string) (result []*models.Person, err error)
	SearchOrganizations(request, since, desc string) (result []*models.Organization, err error)
	SearchVacancies(request, since, desc string) (result []*models.Vacancy, err error)
}
