package pgModels

import baseModels "joblessness/haha/models/base"

type Vacancy struct {
	ID               uint64
	OrganizationID   uint64
	Name             string
	Description      string
	SalaryFrom       int
	SalaryTo         int
	WithTax          bool
	Responsibilities string
	Conditions       string
	Keywords         string
}

func ToPgVacancy(v *baseModels.Vacancy) *Vacancy {
	return &Vacancy{
		ID:               v.ID,
		OrganizationID:   v.Organization.ID,
		Name:             v.Name,
		Description:      v.Description,
		SalaryFrom:       v.SalaryFrom,
		SalaryTo:         v.SalaryTo,
		WithTax:          v.WithTax,
		Responsibilities: v.Responsibilities,
		Conditions:       v.Conditions,
		Keywords:         v.Keywords,
	}
}

func ToBaseVacancy(v *Vacancy, u *User, o *Organization) *baseModels.Vacancy {
	organization := baseModels.VacancyOrganization{
		ID:     u.ID,
		Tag:    u.Tag,
		Email:  u.Email,
		Phone:  u.Phone,
		Avatar: u.Avatar,
		Name:   o.Name,
		Site:   o.Site,
	}

	return &baseModels.Vacancy{
		ID:               v.ID,
		Organization:     organization,
		Name:             v.Name,
		Description:      v.Description,
		SalaryFrom:       v.SalaryFrom,
		SalaryTo:         v.SalaryTo,
		WithTax:          v.WithTax,
		Responsibilities: v.Responsibilities,
		Conditions:       v.Conditions,
		Keywords:         v.Keywords,
	}
}
