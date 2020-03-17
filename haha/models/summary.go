package models

type Education struct {
	Institution string `json:"institution,omitempty"`
	Speciality string `json:"speciality,omitempty"`
	Graduated string `json:"graduated,omitempty"`
	Type string `json:"type,omitempty"`
}

type Experience struct {
	CompanyName string `json:"company_name,omitempty"`
	Role string `json:"role,omitempty"`
	Responsibilities string `json:"responsibilities,omitempty"`
	Start string `json:"start,omitempty"`
	Stop string `json:"stop,omitempty"`
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
	Birthday string `json:"birthday,omitempty"`
}

type Summary struct {
	ID uint64 `json:"id,omitempty"`
	Author Author `json:"author,omitempty"`
	Keywords string `json:"keywords,omitempty"`
	Educations []Education `json:"educations,omitempty"`
	Experiences []Experience `json:"experiences,omitempty"`
}
