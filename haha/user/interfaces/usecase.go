package userInterfaces

import (
	"joblessness/haha/models"
	"mime/multipart"
)

type UserUseCase interface {
	GetPerson(userID uint64) (*models.Person, error)
	ChangePerson(p *models.Person) error
	GetOrganization(userID uint64) (*models.Organization, error)
	ChangeOrganization(o *models.Organization) error
	GetListOfOrgs(page int) (models.Organizations, error)
	SetAvatar(form *multipart.Form, userID uint64) error
	LikeUser(userID, favoriteID uint64) (bool, error)
	LikeExists(userID, favoriteID uint64) (bool, error)
	GetUserFavorite(userID uint64) (models.Favorites, error)
}
