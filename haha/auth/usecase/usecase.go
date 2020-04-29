package authUseCase

import (
	"fmt"
	"google.golang.org/grpc/status"
	"joblessness/haha/auth/interfaces"
	"math/rand"
)

type AuthUseCase struct {
	userRepo authInterfaces.AuthRepository
}

func NewAuthUseCase(authRepo authInterfaces.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: authRepo,
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

func (a *AuthUseCase) RegisterPerson(login, password, name string) (err error) {
	err = a.userRepo.DoesUserExists(login)

	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.AlreadyExists {
		return authInterfaces.ErrUserAlreadyExists
	} else if err != nil {
		return err
	}

	return a.userRepo.RegisterPerson(login, password, name)
}

func (a *AuthUseCase) RegisterOrganization(login, password, name string) (err error) {
	err = a.userRepo.DoesUserExists(login)

	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.AlreadyExists {
		return authInterfaces.ErrUserAlreadyExists
	} else if err != nil {
		return err
	}

	return a.userRepo.RegisterOrganization(login, password, name)
}

func (a *AuthUseCase) Login(login, password string) (userID uint64, role, sessionId string, err error) {
	sessionId = GetSID(64)
	userID, err = a.userRepo.Login(login, password, sessionId)

	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.WrongLoginOrPassword {
		return userID, role, sessionId, authInterfaces.ErrWrongLoginOrPassword
	} else if err == nil {
		role, err = a.userRepo.GetRole(userID)
	}

	return userID, role, sessionId, err
}

func (a *AuthUseCase) Logout(sessionId string) error {
	return a.userRepo.Logout(sessionId)
}

func (a *AuthUseCase) SessionExists(sessionId string) (userID uint64, err error) {
	userID, err = a.userRepo.SessionExists(sessionId)

	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.WrongSID {
		return userID, authInterfaces.ErrWrongSID
	}

	return userID, err
}

func (a *AuthUseCase) GetRole(userID uint64) (role string, err error) {
	role, err = a.userRepo.GetRole(userID)

	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.NotFound {
		return role, authInterfaces.ErrNotFound
	}

	return role, err
}

func (a *AuthUseCase) PersonSession(sessionId string) (uint64, error) {
	userID, err := a.userRepo.SessionExists(sessionId)

	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.WrongSID {
		return userID, authInterfaces.ErrWrongSID
	} else if err != nil {
		return 0, err
	}

	role, err := a.userRepo.GetRole(userID)
	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.NotFound {
		return userID, authInterfaces.ErrNotFound
	} else if err != nil {
		return userID, err
	}

	if role == "person" {
		return userID, nil
	}
	return userID, fmt.Errorf("%w, user id: %d", authInterfaces.ErrUserNotPerson, userID)
}

func (a *AuthUseCase) OrganizationSession(sessionId string) (uint64, error) {
	userID, err := a.userRepo.SessionExists(sessionId)

	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.WrongSID {
		return userID, authInterfaces.ErrWrongSID
	} else if err != nil {
		return 0, err
	}

	role, err := a.userRepo.GetRole(userID)
	if e, ok := status.FromError(err); ok && e.Code() == authInterfaces.NotFound {
		return userID, authInterfaces.ErrNotFound
	} else if err != nil {
		return userID, err
	}

	if role == "organization" {
		return userID, nil
	}
	return userID, fmt.Errorf("%w, user id: %d", authInterfaces.ErrUserNotOrganization, userID)
}
