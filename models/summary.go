package models

type Summary struct {
	UserID uint `json:"author,omitempty"`
	ID uint `json:"id,omitempty"`
	FirstName string `json:"first-name"`
	LastName string `json:"last-name"`
	PhoneNumber string `json:"phone-number"`
	Email string `json:"email"`
	BirthDate string `json:"birth-date"`
	Gender string `json:"gender"`
	Experience string `json:"experience"`
	Education string `json:"education"`
}
