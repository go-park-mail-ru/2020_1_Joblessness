package userInterfaces

import (
	"joblessness/haha/models/base"
)

type UserRepository interface {
	GetPerson(userID uint64) (*baseModels.Person, error)
	ChangePerson(p *baseModels.Person) error
	GetOrganization(userID uint64) (*baseModels.Organization, error)
	ChangeOrganization(p *baseModels.Organization) error
	GetListOfOrgs(page int) (baseModels.Organizations, error)
	SaveAvatarLink(link string, userID uint64) error
	SetOrDeleteLike(userID, favoriteID uint64) (bool, error)
	LikeExists(userID, favoriteID uint64) (bool, error)
	GetUserFavorite(userID uint64) (baseModels.Favorites, error)
}
