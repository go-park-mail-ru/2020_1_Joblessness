package auth

type UseCase interface {
	RegisterPerson(login, password, firstName, lastName, email, phone string) error
	Login(login, password string) (int, string, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (int, error)
}