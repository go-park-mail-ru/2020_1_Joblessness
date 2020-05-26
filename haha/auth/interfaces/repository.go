package authInterfaces

type AuthRepository interface {
	RegisterPerson(login, password, name string) error
	RegisterOrganization(login, password, name string) error
	Login(login, password, SID string) (uint64, error)
	Logout(sessionID string) error
	SessionExists(sessionID string) (uint64, error)
	DoesUserExists(login string) error
	GetRole(userID uint64) (string, error)
}
