package models

type UserPage struct {
	User interface{} `json:"user,omitempty"`
	Summaries interface{} `json:"summaries"`
}

type UserInfo struct {
	Firstname string `json:"firstName,omitempty"`
	Lastname string `json:"lastName,omitempty"`
	Tag string `json:"tag,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type UserSummary struct {
	Title string `json:"title"`
}

type OrganizationInfo struct {
	Name string `json:"name,omitempty"`
	Site string `json:"site,omitempty"`
	Tag string `json:"tag,omitempty"`
	About string `json:"about,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}