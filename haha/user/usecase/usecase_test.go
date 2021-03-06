package userUseCase

//go:generate mockgen -destination=../repository/mock/user.go -package=mock joblessness/haha/user/interfaces UserRepository

import (
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"joblessness/haha/models/base"
	"joblessness/haha/user/repository/mock"
	"mime/multipart"
	"testing"
)

func TestAuthPersonFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	policy := bluemonday.UGCPolicy()
	uc := NewUserUseCase(repo, policy)

	login := "user"
	password := "password"
	phone := "phone"
	firstName := "name"
	userIDEx := uint64(1)
	person := &baseModels.Person{
		Login:     login,
		Password:  password,
		FirstName: firstName,
		Phone:     phone,
	}

	//GetPerson
	repo.EXPECT().GetPerson(userIDEx).Return(person, nil).Times(1)
	resultPerson, err := uc.GetPerson(userIDEx)
	assert.NoError(t, err)
	assert.ObjectsAreEqual(person, resultPerson)

	//ChangePerson
	lastName := "NaNa"
	person.LastName = lastName
	repo.EXPECT().ChangePerson(person).Return(nil).Times(1)
	err = uc.ChangePerson(person)
	assert.NoError(t, err)
}

func TestAuthOrganizationFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	policy := bluemonday.UGCPolicy()
	uc := NewUserUseCase(repo, policy)

	login := "user"
	password := "password"
	phone := "phone"
	name := "name"
	userIDEx := uint64(1)
	organization := &baseModels.Organization{
		Login:    login,
		Password: password,
		Name:     name,
		Phone:    phone,
	}

	//GetOrganization
	repo.EXPECT().GetOrganization(userIDEx).Return(organization, nil).Times(1)
	resultPerson, err := uc.GetOrganization(userIDEx)
	assert.NoError(t, err)
	assert.ObjectsAreEqual(organization, resultPerson)

	//ChangeOrganization
	newName := "NaNa"
	organization.Name = newName
	repo.EXPECT().ChangeOrganization(organization).Return(nil).Times(1)
	err = uc.ChangeOrganization(organization)
	assert.NoError(t, err)
}

func TestSetAvatarNoFile(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	policy := bluemonday.UGCPolicy()
	uc := NewUserUseCase(repo, policy)

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
	policy := bluemonday.UGCPolicy()
	uc := NewUserUseCase(repo, policy)

	repo.EXPECT().GetListOfOrgs(1).Return(baseModels.Organizations{}, nil).Times(1)
	_, err := uc.GetListOfOrgs(1)
	assert.NoError(t, err)
}

func TestLike(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	policy := bluemonday.UGCPolicy()
	uc := NewUserUseCase(repo, policy)

	repo.EXPECT().SetOrDeleteLike(uint64(1), uint64(5)).Return(true, nil).Times(1)
	_, err := uc.LikeUser(uint64(1), uint64(5))
	assert.NoError(t, err)

	repo.EXPECT().GetUserFavorite(uint64(5)).Return(baseModels.Favorites{}, nil).Times(1)
	_, err = uc.GetUserFavorite(uint64(5))
	assert.NoError(t, err)
}
