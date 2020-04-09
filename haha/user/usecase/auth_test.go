package userUseCase

//go:generate mockgen -destination=../repository/mock/user.go -package=mock joblessness/haha/user/interfaces UserRepository

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"joblessness/haha/user/repository/mock"
	"joblessness/haha/models"
	"mime/multipart"
	"testing"
)

func TestAuthPersonFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	uc := NewUserUseCase(repo)

	login := "user"
	password := "password"
	phone := "phone"
	firstName := "name"
	userIdEx := uint64(1)
	person := &models.Person{
		Login:       login,
		Password:    password,
		FirstName: firstName,
		Phone: phone,
	}

	//GetPerson
	repo.EXPECT().GetPerson(userIdEx).Return(person, nil).Times(1)
	resultPerson, err := uc.GetPerson(userIdEx)
	assert.NoError(t, err)
	assert.ObjectsAreEqual(person, resultPerson)

	//ChangePerson
	lastName := "NaNa"
	person.LastName = lastName
	repo.EXPECT().ChangePerson(*person).Return(nil).Times(1)
	err = uc.ChangePerson(*person)
	assert.NoError(t, err)
}

func TestAuthOrganizationFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	uc := NewUserUseCase(repo)

	login := "user"
	password := "password"
	phone := "phone"
	name := "name"
	userIdEx := uint64(1)
	organization := &models.Organization{
		Login:       login,
		Password:    password,
		Name: name,
		Phone: phone,
	}

	//GetOrganization
	repo.EXPECT().GetOrganization(userIdEx).Return(organization, nil).Times(1)
	resultPerson, err := uc.GetOrganization(userIdEx)
	assert.NoError(t, err)
	assert.ObjectsAreEqual(organization, resultPerson)

	//ChangeOrganization
	newName := "NaNa"
	organization.Name = newName
	repo.EXPECT().ChangeOrganization(*organization).Return(nil).Times(1)
	err = uc.ChangeOrganization(*organization)
	assert.NoError(t, err)
}

func TestSetAvatarNoFile(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	uc := NewUserUseCase(repo)

	link := "link"
	form := multipart.Form{
		File: map[string][]*multipart.FileHeader{},
	}
	repo.EXPECT().SaveAvatarLink(link, uint64(1)).Return(nil).Times(0)
	err := uc.SetAvatar(&form, uint64(1))
	assert.Error(t, err)
}

func TestListOrgs(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	uc := NewUserUseCase(repo)

	repo.EXPECT().GetListOfOrgs(1).Return([]models.Organization{}, nil).Times(1)
	_, err := uc.GetListOfOrgs(1)
	assert.NoError(t, err)
}

func TestLike(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	uc := NewUserUseCase(repo)

	repo.EXPECT().SetOrDeleteLike(uint64(1), uint64(5)).Return(true, nil).Times(1)
	_, err := uc.LikeUser(uint64(1), uint64(5))
	assert.NoError(t, err)

	repo.EXPECT().GetUserFavorite(uint64(5)).Return(models.Favorites{}, nil).Times(1)
	_, err = uc.GetUserFavorite(uint64(5))
	assert.NoError(t, err)
}