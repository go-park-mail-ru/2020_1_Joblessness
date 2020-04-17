package pgModels

import (
	"joblessness/haha/models/base"
	"strings"
	"time"
)

type User struct {
	ID             uint64
	Login          string
	Password       string
	OrganizationID uint64
	PersonID       uint64
	Tag            string
	Email          string
	Phone          string
	Registered     time.Time
	Avatar         string
}

type Person struct {
	ID       uint64
	Name     string
	Gender   string
	Birthday time.Time
}

type Organization struct {
	ID    uint64
	Name  string
	Site  string
	About string
}

func ToPgPerson(p *baseModels.Person) (*User, *Person) {
	name := p.FirstName
	if p.LastName != "" {
		name += " " + p.LastName
	}

	return &User{
			ID:         p.ID,
			Login:      p.Login,
			Password:   p.Password,
			Tag:        p.Tag,
			Email:      p.Email,
			Phone:      p.Phone,
			Registered: p.Registered,
			Avatar:     p.Avatar,
		},
		&Person{
			Name:     name,
			Gender:   p.Gender,
			Birthday: p.Birthday,
		}
}

func ToPgOrganization(o *baseModels.Organization) (*User, *Organization) {
	return &User{
			ID:         o.ID,
			Login:      o.Login,
			Password:   o.Password,
			Tag:        o.Tag,
			Email:      o.Email,
			Phone:      o.Phone,
			Registered: o.Registered,
			Avatar:     o.Avatar,
		},
		&Organization{
			Name:  o.Name,
			Site:  o.Site,
			About: o.About,
		}
}

func ToBasePerson(u *User, p *Person) *baseModels.Person {
	var lastName, firstName string
	index := strings.Index(p.Name, " ")
	if index > -1 {
		lastName = p.Name[index+1:]
		firstName = p.Name[:index]
	} else {
		firstName = p.Name
	}

	return &baseModels.Person{
		ID:         u.ID,
		Login:      u.Login,
		Password:   u.Password,
		Tag:        u.Tag,
		Email:      u.Email,
		Phone:      u.Phone,
		Registered: u.Registered,
		Avatar:     u.Avatar,
		FirstName:  firstName,
		LastName:   lastName,
		Gender:     p.Gender,
		Birthday:   p.Birthday,
	}
}

func ToBaseOrganization(u *User, o *Organization) *baseModels.Organization {
	return &baseModels.Organization{
		ID:         u.ID,
		Login:      u.Login,
		Password:   u.Password,
		Tag:        u.Tag,
		About:      o.About,
		Email:      u.Email,
		Phone:      u.Phone,
		Registered: u.Registered,
		Avatar:     u.Avatar,
		Name:       o.Name,
		Site:       o.Site,
	}
}
