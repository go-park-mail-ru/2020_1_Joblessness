package authUseCase

import (
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/models"
	"joblessness/haha/utils/sss"
	"math/rand"
	"mime/multipart"
)

type AuthUseCase struct {
	userRepo authInterfaces.AuthRepository
}

func NewAuthUseCase(userRepo authInterfaces.AuthRepository) *AuthUseCase {
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


func (a *AuthUseCase) Login(login, password string) (userID uint64, role, sessionId string, err error) {
	sessionId = GetSID(64)
	userID, err = a.userRepo.Login(login, password, sessionId)
	if err == nil {
		role, err = a.userRepo.GetRole(userID)
	}

	return userID, role, sessionId, err
}

func (a *AuthUseCase) Logout(sessionId string) error {
	return a.userRepo.Logout(sessionId)
}

func (a *AuthUseCase) SessionExists(sessionId string) (userID uint64, err error) {
	return a.userRepo.SessionExists(sessionId)
}

func (a *AuthUseCase) GetRole(userID uint64) (string, error) {
	return a.userRepo.GetRole(userID)
}

func (a *AuthUseCase) PersonSession(sessionId string) (uint64, error) {
	userID, err := a.userRepo.SessionExists(sessionId)
	if err != nil {
		return 0, err
	}

	role, err := a.userRepo.GetRole(userID)
	if err != nil {
		return userID, err
	}

	if role == "person" {
		return userID, nil
	}
	return userID, authInterfaces.ErrUserNotPerson
}

func (a *AuthUseCase) OrganizationSession(sessionId string) (uint64, error) {
	userID, err := a.userRepo.SessionExists(sessionId)
	if err != nil {
		return 0, err
	}

	role, err := a.userRepo.GetRole(userID)
	if err != nil {
		return userID, err
	}

	if role == "organization" {
		return userID, nil
	}
	return userID, authInterfaces.ErrUserNotOrg
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

func (a *AuthUseCase) GetListOfOrgs(page int) (result []models.Organization, err error) {
	return a.userRepo.GetListOfOrgs(page)
}

func (a *AuthUseCase) SetAvatar(form *multipart.Form, userID uint64) (err error) {
	link, err := sss.UploadAvatar(form, userID)
	if err != nil {
		return err
	}

	return a.userRepo.SaveAvatarLink(link, userID)
}

func (a *AuthUseCase) LikeUser(userID, favoriteID uint64) (bool, error) {
	return a.userRepo.SetOrDeleteLike(userID, favoriteID)
}

func (a *AuthUseCase) LikeExists(userID, favoriteID uint64) (bool, error) {
	return a.userRepo.LikeExists(userID, favoriteID)
}

func (a *AuthUseCase) GetUserFavorite(userID uint64) (models.Favorites, error) {
	return a.userRepo.GetUserFavorite(userID)
}