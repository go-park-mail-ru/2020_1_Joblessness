package authInterfaces

import (
	"joblessness/haha/models"
)

type AuthRepository interface {
	CreatePerson(user *models.Person) error
	CreateOrganization(org *models.Organization) error
	Login(login, password, SID string) (uint64, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	GetRole(userID uint64) (string, error)
	GetPerson(userID uint64) (*models.Person, error)
	ChangePerson(p models.Person) error
	GetOrganization(userID uint64) (*models.Organization, error)
	ChangeOrganization(p models.Organization) error
	DoesUserExists(login string) error
	GetListOfOrgs(page int) ([]models.Organization, error)
	SaveAvatarLink(link string, userID uint64) error
	SetOrDeleteLike(userID, favoriteID uint64) (bool, error)
	GetUserFavorite(userID uint64) (models.Favorites, error)
}
