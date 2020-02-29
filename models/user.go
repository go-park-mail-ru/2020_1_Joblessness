package models

type User struct {
	ID uint `json:"id,omitempty"`
	Login string `json:"login"`
	Password string `json:"password"`

	FirstName string `json:"first-name"`
	LastName string `json:"last-name"`
	Email string `json:"email"`
	PhoneNumber string `json:"phone-number"`
}

type UserLogin struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type UserPage struct {
	User interface{} `json:"user,omitempty"`
	Summaries interface{} `json:"summaries"`
}

type UserInfo struct {
	Firstname string `json:"first-name,omitempty"`
	Lastname string `json:"last-name,omitempty"`
	Tag string `json:"tag,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type UserSummary struct {
	Title string `json:"title"`
}
