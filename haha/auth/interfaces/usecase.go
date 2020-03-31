package authInterfaces

import (
	"joblessness/haha/models"
	"mime/multipart"
)

type AuthUseCase interface {
	RegisterPerson(*models.Person) error
	RegisterOrganization(*models.Organization) error
	Login(login, password string) (uint64, string, error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	PersonSession(sessionId string) (uint64, error)
	OrganizationSession(sessionId string) (uint64, error)
	GetPerson(userID uint64) (*models.Person, error)
	ChangePerson(p models.Person) error
	GetOrganization(userID uint64) (*models.Organization, error)
	ChangeOrganization(o models.Organization) error
	GetListOfOrgs(page int) ([]models.Organization, error)
	SetAvatar(form *multipart.Form, userID uint64) error
	LikeUser(userID, favoriteID uint64) (bool, error)
	GetUserFavorite(userID uint64) (models.Favorites, error)
}