package baseModels

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

//easyjson:json
type Education struct {
	Institution string    `json:"institution,omitempty" validate:"max=60"`
	Speciality  string    `json:"speciality,omitempty" validate:"max=60"`
	Graduated   time.Time `json:"graduated,omitempty"`
	Type        string    `json:"type,omitempty" validate:"max=60"`
}

func (s *Education) Sanitize(policy *bluemonday.Policy) {
	s.Institution = policy.Sanitize(s.Institution)
	s.Speciality = policy.Sanitize(s.Speciality)
	s.Type = policy.Sanitize(s.Type)
}

//easyjson:json
type Experience struct {
	CompanyName      string    `json:"companyName,omitempty" validate:"max=60"`
	Role             string    `json:"role,omitempty" validate:"max=120"`
	Responsibilities string    `json:"responsibilities,omitempty"`
	Start            time.Time `json:"start,omitempty"`
	Stop             time.Time `json:"stop,omitempty"`
}

func (s *Experience) Sanitize(policy *bluemonday.Policy) {
	s.CompanyName = policy.Sanitize(s.CompanyName)
	s.Role = policy.Sanitize(s.Role)
	s.Responsibilities = policy.Sanitize(s.Responsibilities)
}

//easyjson:json
type Author struct {
	ID        uint64    `json:"id,omitempty"`
	Tag       string    `json:"tag,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Gender    string    `json:"gender,omitempty"`
	Birthday  time.Time `json:"birthday,omitempty"`
}

func (s *Author) Sanitize(policy *bluemonday.Policy) {
	s.Tag = policy.Sanitize(s.Tag)
	s.Email = policy.Sanitize(s.Email)
	s.Phone = policy.Sanitize(s.Phone)
	s.Avatar = policy.Sanitize(s.Avatar)
	s.FirstName = policy.Sanitize(s.FirstName)
	s.LastName = policy.Sanitize(s.LastName)
	s.Gender = policy.Sanitize(s.Gender)
}

//easyjson:json
type Summary struct {
	ID          uint64       `json:"id,omitempty"`
	Author      Author       `json:"author,omitempty" validate:"required"`
	Name        string       `json:"name,omitempty" validate:"max=120"`
	SalaryFrom  int          `json:"salaryFrom,omitempty"`
	SalaryTo    int          `json:"salaryTo,omitempty"`
	Keywords    string       `json:"keywords,omitempty"`
	Educations  []Education  `json:"educations,omitempty"`
	Experiences []Experience `json:"experiences,omitempty"`
}

func (s *Summary) Sanitize(policy *bluemonday.Policy) {
	s.Author.Sanitize(policy)
	for _, v := range s.Educations {
		v.Sanitize(policy)
	}
	for _, v := range s.Experiences {
		v.Sanitize(policy)
	}

	s.Name = policy.Sanitize(s.Name)
	s.Keywords = policy.Sanitize(s.Keywords)
}

//easyjson:json
type Summaries []*Summary

func (s *Summaries) Sanitize(policy *bluemonday.Policy) {
	for _, v := range *s {
		v.Sanitize(policy)
	}
}

//easyjson:json
type SendSummary struct {
	VacancyID      uint64    `json:"vacancyId,omitempty"`
	SummaryID      uint64    `json:"summaryId"`
	UserID         uint64    `json:"user_id,omitempty"`
	OrganizationID uint64    `json:"organizationId,omitempty"`
	InterviewDate  time.Time `json:"interview_date,omitempty"`
	Accepted       bool      `json:"accepted"`
	Denied         bool      `json:"denied"`
}

//easyjson:json
type VacancyResponse struct {
	UserID      uint64 `json:"user_id,omitempty"`
	Tag         string `json:"tag,omitempty"`
	VacancyID   uint64 `json:"vacancyId,omitempty"`
	SummaryID   uint64 `json:"summaryId,omitempty"`
	FirstName    string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Accepted    bool   `json:"accepted"`
	Denied      bool   `json:"denied"`
}

func (s *VacancyResponse) Sanitize(policy *bluemonday.Policy) {
	s.Tag = policy.Sanitize(s.Tag)
	s.FirstName = policy.Sanitize(s.FirstName)
	s.LastName = policy.Sanitize(s.LastName)
	s.Avatar = policy.Sanitize(s.Avatar)
}

//easyjson:json
type OrgSummaries []*VacancyResponse

func (s *OrgSummaries) Sanitize(policy *bluemonday.Policy) {
	for _, v := range *s {
		v.Sanitize(policy)
	}
}
