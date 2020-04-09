package authInterfaces

import "joblessness/haha/models"

type AuthRepository interface {
	CreatePerson(user *models.Person) error
	CreateOrganization(org *models.Organization) error
	Login(login, password, SID string) (uint64, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	DoesUserExists(login string) error
	GetRole(userID uint64) (string, error)
}
