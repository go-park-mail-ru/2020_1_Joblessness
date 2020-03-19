package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"joblessness/haha/auth/repository/mock"
	"joblessness/haha/models"
	"mime/multipart"
	"testing"
)

func TestAuthPersonFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	uc := NewAuthUseCase(repo)

	login := "user"
	password := "password"
	phone := "phone"
	firstName := "name"
	sidEx := "sid"
	userIdEx := uint64(1)
	person := &models.Person{
		Login:       login,
		Password:    password,
		FirstName: firstName,
		Phone: phone,
	}

	//RegisterPerson
	repo.EXPECT().CreatePerson(person).Return(nil).Times(1)
	repo.EXPECT().DoesUserExists(login).Return(nil).Times(1)
	err := uc.RegisterPerson(person)
	assert.NoError(t, err)

	//Login
	repo.EXPECT().Login(login, password, gomock.Any()).Return(userIdEx, nil).Times(1)
	userId, sid, err := uc.Login(login, password)
	assert.NoError(t, err)
	assert.Equal(t, userIdEx, userId, "Id corrupted")
	assert.NotEmpty(t, sid, "No sid")

	//Logout
	repo.EXPECT().Logout(sidEx).Return(nil).Times(1)
	err = uc.Logout(sidEx)
	assert.NoError(t, err)

	//Check
	repo.EXPECT().SessionExists(sidEx).Return(userIdEx, nil).Times(1)
	userId, err = uc.SessionExists(sidEx)
	assert.NoError(t, err)
	assert.Equal(t, userIdEx, userId, "Id corrupted")

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
	uc := NewAuthUseCase(repo)

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

	//RegisterOrganization
	repo.EXPECT().CreateOrganization(organization).Return(nil).Times(1)
	repo.EXPECT().DoesUserExists(login).Return(nil).Times(1)
	err := uc.RegisterOrganization(organization)
	assert.NoError(t, err)

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
	uc := NewAuthUseCase(repo)

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
	uc := NewAuthUseCase(repo)

	repo.EXPECT().GetListOfOrgs(1).Return([]models.Organization{}, nil).Times(1)
	_, err := uc.GetListOfOrgs(1)
	assert.NoError(t, err)
}