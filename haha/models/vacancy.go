package models

type VacancyOrganization struct {
	ID uint64 `json:"id,omitempty"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Name string `json:"name,omitempty"`
	Site string `json:"site,omitempty"`
}

type Vacancy struct {
	ID uint64 `json:"id,omitempty"`
	Organization VacancyOrganization `json:"organization,omitempty"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SalaryFrom int `json:"salary_from,omitempty"`
	SalaryTo int `json:"salary_to,omitempty"`
	WithTax bool `json:"with_tax,omitempty"`
	Responsibilities string `json:"responsibilities,omitempty"`
	Conditions string `json:"conditions,omitempty"`
	Keywords string`json:"keywords,omitempty"`
}
