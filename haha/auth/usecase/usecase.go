package usecase

import (
	"joblessness/haha/auth"
	"joblessness/haha/models"
	"math/rand"
)

type AuthUseCase struct {
	userRepo auth.UserRepository
}

func NewAuthUseCase(userRepo auth.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo:userRepo,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetSID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (a *AuthUseCase) RegisterPerson(login, password, firstName, lastName, email, phone string) (err error) {
	person := &models.Person{
		ID:          0,
		Login:       login,
		Password:    password,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phone,
	}

	return a.userRepo.CreatePerson(person)
}

func (a *AuthUseCase) Login(login, password string) (userId int, sessionId string, err error) {
	sessionId = GetSID(64)
	userId, err = a.userRepo.Login(login, password, sessionId)
	return userId, sessionId, err
}

func (a *AuthUseCase) Logout(sessionId string) error {
	return a.userRepo.Logout(sessionId)
}

func (a *AuthUseCase) SessionExists(sessionId string) (int, error) {
	return a.userRepo.SessionExists(sessionId)
}