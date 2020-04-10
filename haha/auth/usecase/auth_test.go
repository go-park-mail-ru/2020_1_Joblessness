package authUseCase

//go:generate mockgen -destination=../repository/mock/user.go -package=mock joblessness/haha/user/interfaces UserRepository

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"joblessness/haha/auth/repository/mock"
	"joblessness/haha/models"
	"testing"
)

func TestAuthPersonFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockAuthRepository(controller)
	uc := NewAuthUseCase(repo)

	login := "user"
	password := "password"
	phone := "phone"
	firstName := "name"
	sidEx := "sid"
	userIdEx := uint64(1)
	person := &models.Person{
		Login:     login,
		Password:  password,
		FirstName: firstName,
		Phone:     phone,
	}

	//RegisterPerson
	repo.EXPECT().CreatePerson(person).Return(nil).Times(1)
	repo.EXPECT().DoesUserExists(login).Return(nil).Times(1)
	err := uc.RegisterPerson(person)
	assert.NoError(t, err)

	//Login
	repo.EXPECT().Login(login, password, gomock.Any()).Return(userIdEx, nil).Times(1)
	repo.EXPECT().GetRole(userIdEx).Return("userIdEx", nil).Times(1)
	userId, _, sid, err := uc.Login(login, password)
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
}

func TestAuthOrganizationFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockAuthRepository(controller)
	uc := NewAuthUseCase(repo)

	login := "user"
	password := "password"
	phone := "phone"
	name := "name"
	organization := &models.Organization{
		Login:    login,
		Password: password,
		Name:     name,
		Phone:    phone,
	}

	//RegisterOrganization
	repo.EXPECT().CreateOrganization(organization).Return(nil).Times(1)
	repo.EXPECT().DoesUserExists(login).Return(nil).Times(1)
	err := uc.RegisterOrganization(organization)
	assert.NoError(t, err)
}
