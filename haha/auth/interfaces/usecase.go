package authInterfaces

type AuthUseCase interface {
	RegisterPerson(login, password, name string) error
	RegisterOrganization(login, password, name string) error
	Login(login, password string) (userID uint64, role, sessionID string, err error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	PersonSession(sessionId string) (uint64, error)
	OrganizationSession(sessionId string) (uint64, error)
	GetRole(userID uint64) (string, error)
}
