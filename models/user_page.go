package models

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
