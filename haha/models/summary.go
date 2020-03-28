package models

import "time"

type Education struct {
	Institution string `json:"institution,omitempty"`
	Speciality  string `json:"speciality,omitempty"`
	Graduated   time.Time   `json:"graduated,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Experience struct {
	CompanyName string `json:"company_name,omitempty"`
	Role string `json:"role,omitempty"`
	Responsibilities string `json:"responsibilities,omitempty"`
	Start time.Time `json:"start,omitempty"`
	Stop time.Time `json:"stop,omitempty"`
}

type Author struct {
	ID uint64 `json:"id,omitempty"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
}

type Summary struct {
	ID uint64 `json:"id,omitempty"`
	Author Author `json:"author,omitempty"`
	Keywords string `json:"keywords,omitempty"`
	Educations []Education `json:"educations,omitempty"`
	Experiences []Experience `json:"experiences,omitempty"`
}

type SendSummary struct {
	VacancyID uint64 `json:"vacancy_id"`
	SummaryID uint64 `json:"summary_id"`
	UserID uint64 `json:"user_id,omitempty"`
	OrganizationID uint64 `json:"organization_id,omitempty"`
	Accepted bool `json:"accepted,omitempty"`
	Denied bool `json:"denied,omitempty"`
}

type VacancyResponse struct {
	UserID uint64 `json:"user_id,omitempty"`
	Tag string `json:"tag,omitempty"`
	VacancyID uint64 `json:"vacancy_id"`
	SummaryID uint64 `json:"summary_id"`

	Keywords  string `json:"keywords,omitempty"`
}

type OrgSummaries []*VacancyResponse