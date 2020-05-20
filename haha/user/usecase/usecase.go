package userUseCase

import (
	"github.com/microcosm-cc/bluemonday"
	"joblessness/haha/models/base"
	"joblessness/haha/user/interfaces"
	"joblessness/haha/utils/sss"
	"mime/multipart"
)

type UserUseCase struct {
	userRepo userInterfaces.UserRepository
	policy   *bluemonday.Policy
}

func NewUserUseCase(userRepo userInterfaces.UserRepository, policy *bluemonday.Policy) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		policy:   policy,
	}
}

func (u *UserUseCase) GetPerson(userID uint64) (*baseModels.Person, error) {
	res, err := u.userRepo.GetPerson(userID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *UserUseCase) ChangePerson(p *baseModels.Person) error {
	return u.userRepo.ChangePerson(p)
}

func (u *UserUseCase) GetOrganization(userID uint64) (*baseModels.Organization, error) {
	res, err := u.userRepo.GetOrganization(userID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *UserUseCase) ChangeOrganization(o *baseModels.Organization) error {
	return u.userRepo.ChangeOrganization(o)
}

func (u *UserUseCase) GetListOfOrgs(page int) (result baseModels.Organizations, err error) {
	res, err := u.userRepo.GetListOfOrgs(page)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *UserUseCase) SetAvatar(form *multipart.Form, userID uint64) (err error) {
	link, err := sss.UploadAvatar(form, userID)
	if err != nil {
		return err
	}

	return u.userRepo.SaveAvatarLink(link, userID)
}

func (u *UserUseCase) LikeUser(userID, favoriteID uint64) (bool, error) {
	return u.userRepo.SetOrDeleteLike(userID, favoriteID)
}

func (u *UserUseCase) LikeExists(userID, favoriteID uint64) (bool, error) {
	return u.userRepo.LikeExists(userID, favoriteID)
}

func (u *UserUseCase) GetUserFavorite(userID uint64) (baseModels.Favorites, error) {
	return u.userRepo.GetUserFavorite(userID)
}
