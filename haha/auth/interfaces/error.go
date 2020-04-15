package authInterfaces

import (
	"errors"
)

var (
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrWrongSID = errors.New("wrong sid")
	ErrUserNotPerson = errors.New("user is not a person")
	ErrUserNotOrganization = errors.New("user is not a organization")
	ErrUserAlreadyExists = errors.New("user already exists")
)
