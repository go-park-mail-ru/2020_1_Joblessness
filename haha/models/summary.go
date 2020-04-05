package models

import "time"

type Education struct {
	Institution string `json:"institution,omitempty"`
	Speciality  string `json:"speciality,omitempty"`
	Graduated   time.Time   `json:"graduated,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Experience struct {
	CompanyName string `json:"companyName,omitempty"`
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
	FirstName string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	Gender string `json:"gender,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
}

type Summary struct {
	ID uint64 `json:"id,omitempty"`
	Author Author `json:"author,omitempty" validate:"required"`
	Keywords string `json:"keywords,omitempty"`
	Educations []Education `json:"educations,omitempty"`
	Experiences []Experience `json:"experiences,omitempty"`
}

type SendSummary struct {
	VacancyID uint64 `json:"vacancyId,omitempty"`
	SummaryID uint64 `json:"summaryId"`
	UserID uint64 `json:"user_id,omitempty"`
	OrganizationID uint64 `json:"organizationId,omitempty"`
	Accepted bool `json:"accepted"`
	Denied bool `json:"denied"`
}

type VacancyResponse struct {
	UserID uint64 `json:"user_id,omitempty"`
	Tag string `json:"tag,omitempty"`
	VacancyID uint64 `json:"vacancyId,omitempty"`
	SummaryID uint64 `json:"summaryId,omitempty"`
	Keywords  string `json:"keywords,omitempty"`
	VacancyName  string `json:"vacancyName,omitempty"`
	SummaryName  string `json:"summaryName,omitempty"`
}

type OrgSummaries []*VacancyResponse