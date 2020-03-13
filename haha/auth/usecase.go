package auth

import "joblessness/haha/models"

type UseCase interface {
	RegisterPerson(login, password, firstName, lastName, email, phone string) error
	Login(login, password string) (uint64, string, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	GetPerson(userID uint64) (models.Person, error)
	ChangePerson(p models.Person) error
}