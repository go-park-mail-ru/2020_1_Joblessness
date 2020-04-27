package baseModels

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

//easyjson:json
type Person struct {
	ID         uint64    `json:"id,omitempty"`
	Login      string    `json:"login,omitempty" validate:"required,max=60"`
	Password   string    `json:"password,omitempty" validate:"required,max=60"`
	Tag        string    `json:"tag,omitempty" validate:"omitempty,min=6,max=20"`
	Email      string    `json:"email,omitempty" validate:"omitempty,max=60,email"`
	Phone      string    `json:"phone,omitempty" validate:"max=20"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
	FirstName  string    `json:"firstName,omitempty" validate:"max=60"`
	LastName   string    `json:"lastName,omitempty" validate:"max=60"`
	Gender     string    `json:"gender,omitempty" validate:"max=10"`
	Birthday   time.Time `json:"birthday,omitempty"`
}

func (s *Person) Sanitize(policy *bluemonday.Policy) {
	s.Login = policy.Sanitize(s.Login)
	s.Tag = policy.Sanitize(s.Tag)
	s.Email = policy.Sanitize(s.Email)
	s.Phone = policy.Sanitize(s.Phone)
	s.FirstName = policy.Sanitize(s.FirstName)
	s.LastName = policy.Sanitize(s.LastName)
	s.Gender = policy.Sanitize(s.Gender)
}

//easyjson:json
type Persons []*Person

func (s *Persons) Sanitize(policy *bluemonday.Policy) {
	for _, v := range *s {
		v.Sanitize(policy)
	}
}

//easyjson:json
type Organization struct {
	ID         uint64    `json:"id,omitempty"`
	Login      string    `json:"login,omitempty" validate:"required,max=60"`
	Password   string    `json:"password,omitempty" validate:"required,max=60"`
	Tag        string    `json:"tag,omitempty" validate:"omitempty,min=6,max=20"`
	Email      string    `json:"email,omitempty" validate:"max=60"`
	Phone      string    `json:"phone,omitempty" validate:"max=20"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
	Name       string    `json:"name,omitempty"`
	About      string    `json:"about,omitempty"`
	Site       string    `json:"site,omitempty"`
}

func (s *Organization) Sanitize(policy *bluemonday.Policy) {
	s.Login = policy.Sanitize(s.Login)
	s.Tag = policy.Sanitize(s.Tag)
	s.Email = policy.Sanitize(s.Email)
	s.Phone = policy.Sanitize(s.Phone)
	s.Name = policy.Sanitize(s.Name)
	s.About = policy.Sanitize(s.About)
	s.Site = policy.Sanitize(s.Site)
}

//easyjson:json
type Organizations []*Organization

func (s *Organizations) Sanitize(policy *bluemonday.Policy) {
	for _, v := range *s {
		v.Sanitize(policy)
	}
}

//easyjson:json
type UserLogin struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

//easyjson:json
type Favorite struct {
	ID       uint64 `json:"id"`
	Tag      string `json:"tag"`
	Avatar string `json:"avatar"`
	IsPerson bool   `json:"isPerson"`
	Name string `json:"name"`
	Surname string `json:"surname"`
}

//easyjson:json
type Favorites []*Favorite

//easyjson:json
type ResponseRole struct {
	ID   uint64 `json:"id"`
	Role string `json:"role"`
}

type Role struct {
	Person       bool
	Organization bool
}
