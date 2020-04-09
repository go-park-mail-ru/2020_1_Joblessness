package userUseCase

import (
	"joblessness/haha/models"
	"joblessness/haha/user/interfaces"
	"joblessness/haha/utils/sss"
	"mime/multipart"
)

type UserUseCase struct {
	userRepo userInterfaces.UserRepository
}

func NewUserUseCase(userRepo userInterfaces.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo:userRepo,
	}
}

func (a *UserUseCase) GetPerson(userID uint64) (*models.Person, error) {
	return a.userRepo.GetPerson(userID)
}

func (a *UserUseCase) ChangePerson(p models.Person) error {
	return a.userRepo.ChangePerson(p)
}

func (a *UserUseCase) GetOrganization(userID uint64) (*models.Organization, error) {
	return a.userRepo.GetOrganization(userID)
}

func (a *UserUseCase) ChangeOrganization(o models.Organization) error {
	return a.userRepo.ChangeOrganization(o)
}

func (a *UserUseCase) GetListOfOrgs(page int) (result []models.Organization, err error) {
	return a.userRepo.GetListOfOrgs(page)
}

func (a *UserUseCase) SetAvatar(form *multipart.Form, userID uint64) (err error) {
	link, err := sss.UploadAvatar(form, userID)
	if err != nil {
		return err
	}

	return a.userRepo.SaveAvatarLink(link, userID)
}

func (a *UserUseCase) LikeUser(userID, favoriteID uint64) (bool, error) {
	return a.userRepo.SetOrDeleteLike(userID, favoriteID)
}

func (a *UserUseCase) LikeExists(userID, favoriteID uint64) (bool, error) {
	return a.userRepo.LikeExists(userID, favoriteID)
}

func (a *UserUseCase) GetUserFavorite(userID uint64) (models.Favorites, error) {
	return a.userRepo.GetUserFavorite(userID)
}
