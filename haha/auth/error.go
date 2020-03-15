package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrWrongSID    = errors.New("wrong sid")
	ErrWrongLogPas      = errors.New("wrong login or password")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotPerson = errors.New("this user is not a person")
	ErrUserNotOrg = errors.New("this user is not a organization")
)