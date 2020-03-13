package models

type Vacancy struct {
	ID uint64 `json:"id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description"`
	Skills string `json:"skills"`
	Salary int `json:"salary"`
	Address string `json:"address"`
	PhoneNumber string `json:"phone-number"`
}
