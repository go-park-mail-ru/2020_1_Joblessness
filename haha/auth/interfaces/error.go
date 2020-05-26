package authInterfaces

import (
	"errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrWrongSID             = errors.New("wrong sid")
	ErrUserNotPerson        = errors.New("user is not a person")
	ErrUserNotOrganization  = errors.New("user is not a organization")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrParseGrpcError       = errors.New("can't parse rpc error")
	ErrNotFound             = errors.New("user not found")

	AlreadyExists        codes.Code = 400
	WrongLoginOrPassword codes.Code = 400
	NotFound             codes.Code = 404
	WrongSID             codes.Code = 500
)
