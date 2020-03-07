package auth

import (
	"joblessness/haha/models"
)

type UserRepository interface {
	CreatePerson(user *models.Person) error
	Login(login, password, SID string) (int, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (int, error)
}
