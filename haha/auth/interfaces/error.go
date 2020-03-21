package authInterfaces

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrWrongSID    = errors.New("wrong sid")
	ErrWrongLogPas      = errors.New("wrong login or password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotPerson = errors.New("this user is not a person")
	ErrUserNotOrg = errors.New("this user is not a organization")
	ErrAvatarParse = errors.New("not able to parse avatar")
	ErrUploadAvatar = errors.New("not able to upload avatar")
	ErrNoFile = errors.New("no file in multipart form")
)