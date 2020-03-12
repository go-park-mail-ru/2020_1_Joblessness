package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"joblessness/haha/auth/repository/mock"
	"joblessness/haha/models"
	"testing"
)

func TestAuthFlow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := mock.NewMockUserRepository(controller)
	uc := NewAuthUseCase(repo)

	login := "user"
	password := "password"
	sidEx := "sid"
	userIdEx := 1
	person := &models.Person{
		Login:       login,
		Password:    password,
	}

	//Register
	repo.EXPECT().CreatePerson(person).Return(nil).Times(1)
	err := uc.RegisterPerson(login, password, "", "", "", "")
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
}