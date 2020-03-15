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

func (a *AuthUseCase) RegisterPerson(p *models.Person) (err error) {
	err = a.userRepo.DoesUserExists(p.Login)
	if err != nil {
		return err
	}

	return a.userRepo.CreatePerson(p)
}

func (a *AuthUseCase) RegisterOrganization(o *models.Organization) (err error) {
	err = a.userRepo.DoesUserExists(o.Login)
	if err != nil {
		return err
	}

	return a.userRepo.CreateOrganization(o)
}


func (a *AuthUseCase) Login(login, password string) (userId uint64, sessionId string, err error) {
	sessionId = GetSID(64)
	userId, err = a.userRepo.Login(login, password, sessionId)
	return userId, sessionId, err
}

func (a *AuthUseCase) Logout(sessionId string) error {
	return a.userRepo.Logout(sessionId)
}

func (a *AuthUseCase) SessionExists(sessionId string) (uint64, error) {
	return a.userRepo.SessionExists(sessionId)
}

func (a *AuthUseCase) GetPerson(userID uint64) (*models.Person, error) {
	return a.userRepo.GetPerson(userID)
}

func (a *AuthUseCase) ChangePerson(p models.Person) error {
	return a.userRepo.ChangePerson(p)
}

func (a *AuthUseCase) GetOrganization(userID uint64) (*models.Organization, error) {
	return a.userRepo.GetOrganization(userID)
}

func (a *AuthUseCase) ChangeOrganization(o models.Organization) error {
	return a.userRepo.ChangeOrganization(o)
}