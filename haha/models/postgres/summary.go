package pgModels

import (
	"database/sql"
	"joblessness/haha/models/base"
	"strings"
)

type Summary struct {
	ID         uint64
	AuthorID   uint64
	Keywords   string
	Name       string
	SalaryFrom int
	SalaryTo   int
}

type Education struct {
	SummaryID   uint64
	Institution string
	Speciality  string
	Graduated   sql.NullTime
	Type        string
}

type Experience struct {
	SummaryID        uint64
	CompanyName      string
	Role             string
	Responsibilities string
	Start            sql.NullTime
	Stop             sql.NullTime
}

func ToPgSummary(s *baseModels.Summary) (summary *Summary, educations []*Education, experiences []*Experience) {
	summary = &Summary{
		ID:         s.ID,
		AuthorID:   s.Author.ID,
		Keywords:   s.Keywords,
		Name:       s.Name,
		SalaryFrom: s.SalaryFrom,
		SalaryTo:   s.SalaryTo,
	}

	for _, education := range s.Educations {
		educations = append(educations, &Education{
			SummaryID:   summary.ID,
			Institution: education.Institution,
			Speciality:  education.Speciality,
			Graduated:   sql.NullTime{education.Graduated, !education.Graduated.IsZero()},
			Type:        education.Type,
		})
	}

	for _, experience := range s.Experiences {
		experiences = append(experiences, &Experience{
			SummaryID:        summary.ID,
			CompanyName:      experience.CompanyName,
			Role:             experience.Role,
			Responsibilities: experience.Responsibilities,
			Start:            sql.NullTime{experience.Start, !experience.Start.IsZero()},
			Stop:             sql.NullTime{experience.Stop, !experience.Stop.IsZero()},
		})
	}

	return summary, educations, experiences
}

func ToBaseSummary(s *Summary, eds []*Education, exs []*Experience, u *User, p *Person) *baseModels.Summary {
	var educations []baseModels.Education

	for _, ed := range eds {
		educations = append(educations, baseModels.Education{
			Institution: ed.Institution,
			Speciality:  ed.Speciality,
			Graduated:   ed.Graduated.Time,
			Type:        ed.Type,
		})
	}

	var experiences []baseModels.Experience

	for _, ex := range exs {
		experiences = append(experiences, baseModels.Experience{
			CompanyName:      ex.CompanyName,
			Role:             ex.Role,
			Responsibilities: ex.Responsibilities,
			Start:            ex.Start.Time,
			Stop:             ex.Stop.Time,
		})
	}

	nameArr := strings.Split(p.Name, " ")
	firstName := nameArr[0]
	var lastName string
	if len(nameArr) > 1 {
		lastName = nameArr[1]
	}

	author := baseModels.Author{
		ID:        u.ID,
		Tag:       u.Tag,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    p.Gender,
		Birthday:  p.Birthday,
	}

	return &baseModels.Summary{
		ID:          s.ID,
		Author:      author,
		Keywords:    s.Keywords,
		Name:        s.Name,
		SalaryFrom:  s.SalaryFrom,
		SalaryTo:    s.SalaryTo,
		Educations:  educations,
		Experiences: experiences,
	}
}
