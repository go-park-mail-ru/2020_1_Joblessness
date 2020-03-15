package auth

import (
	"joblessness/haha/models"
)

type UserRepository interface {
	CreatePerson(user *models.Person) error
	CreateOrganization(org *models.Organization) error
	Login(login, password, SID string) (uint64, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	GetPerson(userID uint64) (*models.Person, error)
	ChangePerson(p models.Person) error
	GetOrganization(userID uint64) (*models.Organization, error)
	ChangeOrganization(p models.Organization) error
	DoesUserExists(login string) error
}
