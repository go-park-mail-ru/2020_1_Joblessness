package userInterfaces

import (
	"joblessness/haha/models"
)

type UserRepository interface {
	GetPerson(userID uint64) (*models.Person, error)
	ChangePerson(p models.Person) error
	GetOrganization(userID uint64) (*models.Organization, error)
	ChangeOrganization(p models.Organization) error
	GetListOfOrgs(page int) ([]models.Organization, error)
	SaveAvatarLink(link string, userID uint64) error
	SetOrDeleteLike(userID, favoriteID uint64) (bool, error)
	LikeExists(userID, favoriteID uint64) (bool, error)
	GetUserFavorite(userID uint64) (models.Favorites, error)
}
