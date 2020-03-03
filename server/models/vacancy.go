package models

type Vacancy struct {
	ID uint `json:"id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description"`
	Skills string `json:"skills"`
	Salary string `json:"salary"`
	Address string `json:"address"`
	PhoneNumber string `json:"phone-number"`
}
