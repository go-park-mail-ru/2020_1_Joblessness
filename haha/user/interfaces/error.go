package userInterfaces

import "fmt"

type ErrorUserNotFound struct {
	ID uint64
}

func NewErrorUserNotFound(id uint64) *ErrorUserNotFound {
	return &ErrorUserNotFound{ID: id}
}

func (e *ErrorUserNotFound) Error() string {
	return fmt.Sprintf("User with id %d not found", e.ID)
}

type ErrorWrongSID struct {}

func NewErrorWrongSID() *ErrorWrongSID {
	return &ErrorWrongSID{}
}

func (e *ErrorWrongSID) Error() string {
	return "Wrong SID"
}

type ErrorWrongLoginOrPassword struct {}

func NewErrorWrongLoginOrPassword() *ErrorWrongLoginOrPassword {
	return &ErrorWrongLoginOrPassword{}
}

func (e *ErrorWrongLoginOrPassword) Error() string {
	return "Wrong login or password"
}

type ErrorUserAlreadyExists struct {
	Login string
}

func NewErrorUserAlreadyExists(login string) *ErrorUserAlreadyExists {
	return &ErrorUserAlreadyExists{Login: login}
}

func (e *ErrorUserAlreadyExists) Error() string {
	return fmt.Sprintf("User with login %s already exists", e.Login)
}

type ErrorUserNotPerson struct {
	ID uint64
}

func NewErrorUserNotPerson(id uint64) *ErrorUserNotPerson {
	return &ErrorUserNotPerson{ID: id}
}

func (e *ErrorUserNotPerson) Error() string {
	return fmt.Sprintf("User with id %d is not a person", e.ID)
}

type ErrorUserNotOrganization struct {
	ID uint64
}

func NewErrorUserNotOrganization(id uint64) *ErrorUserNotOrganization {
	return &ErrorUserNotOrganization{ID: id}
}

func (e *ErrorUserNotOrganization) Error() string {
	return fmt.Sprintf("User with id %d is not a organization", e.ID)
}

type ErrorAvatarParse struct {}

func NewErrorAvatarParse() *ErrorAvatarParse {
	return &ErrorAvatarParse{}
}

func (e *ErrorAvatarParse) Error() string {
	return "Not able to parse avatar"
}

type ErrorUploadAvatar struct {}

func NewErrorUploadAvatar() *ErrorUploadAvatar {
	return &ErrorUploadAvatar{}
}

func (e *ErrorUploadAvatar) Error() string {
	return "Not able to upload avatar"
}

type ErrorNoFile struct {}

func NewErrorNoFile() *ErrorNoFile {
	return &ErrorNoFile{}
}

func (e *ErrorNoFile) Error() string {
	return "No file in multipart form"
}
