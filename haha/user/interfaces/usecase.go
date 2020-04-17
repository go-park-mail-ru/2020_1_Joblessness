package userInterfaces

import (
	"joblessness/haha/models/base"
	"mime/multipart"
)

type UserUseCase interface {
	GetPerson(userID uint64) (*baseModels.Person, error)
	ChangePerson(p *baseModels.Person) error
	GetOrganization(userID uint64) (*baseModels.Organization, error)
	ChangeOrganization(o *baseModels.Organization) error
	GetListOfOrgs(page int) (baseModels.Organizations, error)
	SetAvatar(form *multipart.Form, userID uint64) error
	LikeUser(userID, favoriteID uint64) (bool, error)
	LikeExists(userID, favoriteID uint64) (bool, error)
	GetUserFavorite(userID uint64) (baseModels.Favorites, error)
}
