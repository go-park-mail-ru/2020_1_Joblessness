package authInterfaces

type AuthUseCase interface {
	RegisterPerson(login, password, name string) error
	RegisterOrganization(login, password, name string) error
	Login(login, password string) (userID uint64, role, sessionID string, err error)
	Logout(sessionID string) error
	SessionExists(sessionID string) (uint64, error)
	PersonSession(sessionID string) (uint64, error)
	OrganizationSession(sessionID string) (uint64, error)
	GetRole(userID uint64) (string, error)
}
