package baseModels

import "github.com/microcosm-cc/bluemonday"

type VacancyOrganization struct {
	ID     uint64 `json:"id,omitempty"`
	Tag    string `json:"tag,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Name   string `json:"name,omitempty" validate:"max=60"`
	Site   string `json:"site,omitempty" validate:"max=60"`
}

func (s *VacancyOrganization) Sanitize(policy *bluemonday.Policy) {
	s.Tag = policy.Sanitize(s.Tag)
	s.Email = policy.Sanitize(s.Email)
	s.Phone = policy.Sanitize(s.Phone)
	s.Name = policy.Sanitize(s.Name)
	s.Site = policy.Sanitize(s.Site)
}

type Vacancy struct {
	ID               uint64              `json:"id,omitempty"`
	Organization     VacancyOrganization `json:"organization,omitempty"`
	Name             string              `json:"name,omitempty" validate:"required,max=60"`
	Description      string              `json:"description,omitempty"`
	SalaryFrom       int                 `json:"salaryFrom,omitempty"`
	SalaryTo         int                 `json:"salaryTo,omitempty"`
	WithTax          bool                `json:"withTax"`
	Responsibilities string              `json:"responsibilities,omitempty"`
	Conditions       string              `json:"conditions,omitempty"`
	Keywords         string              `json:"keywords,omitempty"`
}

func (s *Vacancy) Sanitize(policy *bluemonday.Policy) {
	s.Organization.Sanitize(policy)

	s.Name = policy.Sanitize(s.Name)
	s.Description = policy.Sanitize(s.Description)
	s.Responsibilities = policy.Sanitize(s.Responsibilities)
	s.Conditions = policy.Sanitize(s.Conditions)
	s.Keywords = policy.Sanitize(s.Keywords)
}

type Vacancies []*Vacancy

func (s *Vacancies) Sanitize(policy *bluemonday.Policy) {
	for _, v := range *s {
		v.Sanitize(policy)
	}
}
