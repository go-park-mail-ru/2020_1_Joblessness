package authInterfaces

import (
	"joblessness/haha/models/base"
)

type AuthUseCase interface {
	RegisterPerson(*baseModels.Person) error
	RegisterOrganization(*baseModels.Organization) error
	Login(login, password string) (userID uint64, role, sessionID string, err error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	PersonSession(sessionId string) (uint64, error)
	OrganizationSession(sessionId string) (uint64, error)
	GetRole(userID uint64) (string, error)
}
