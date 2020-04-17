package authInterfaces

import (
	"joblessness/haha/models/base"
)

type AuthRepository interface {
	CreatePerson(user *baseModels.Person) error
	CreateOrganization(org *baseModels.Organization) error
	Login(login, password, SID string) (uint64, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	DoesUserExists(login string) error
	GetRole(userID uint64) (string, error)
}
