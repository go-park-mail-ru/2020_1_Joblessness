package authInterfaces

import (
	"joblessness/haha/models"
	"mime/multipart"
)

type AuthUseCase interface {
	RegisterPerson(*models.Person) error
	RegisterOrganization(*models.Organization) error
	Login(login, password string) (userID uint64, role, sessionID string, err error)
	Logout(sessionId string) error
	SessionExists(sessionId string) (uint64, error)
	PersonSession(sessionId string) (uint64, error)
	OrganizationSession(sessionId string) (uint64, error)
	GetRole(userID uint64) (string, error)
	GetPerson(userID uint64) (*models.Person, error)
	ChangePerson(p models.Person) error
	GetOrganization(userID uint64) (*models.Organization, error)
	ChangeOrganization(o models.Organization) error
	GetListOfOrgs(page int) ([]models.Organization, error)
	SetAvatar(form *multipart.Form, userID uint64) error
<<<<<<< HEAD:haha/auth/usecase.go
	HasPersons() (has bool, err error)
	HasOrganizations() (has bool, err error)
=======
	LikeUser(userID, favoriteID uint64) (bool, error)
	GetUserFavorite(userID uint64) (models.Favorites, error)
>>>>>>> 3f6088f97fc3b00a828300b9b6df976592519559:haha/auth/interfaces/usecase.go
}