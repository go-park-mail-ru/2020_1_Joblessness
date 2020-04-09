package authInterfaces

import "fmt"

type ErrorWrongLoginOrPassword struct {}

func NewErrorWrongLoginOrPassword() *ErrorWrongLoginOrPassword {
	return &ErrorWrongLoginOrPassword{}
}

func (e *ErrorWrongLoginOrPassword) Error() string {
	return "Wrong login or password"
}

type ErrorWrongSID struct {}

func NewErrorWrongSID() *ErrorWrongSID {
	return &ErrorWrongSID{}
}

func (e *ErrorWrongSID) Error() string {
	return "Wrong SID"
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

type ErrorUserAlreadyExists struct {
	Login string
}

func NewErrorUserAlreadyExists(login string) *ErrorUserAlreadyExists {
	return &ErrorUserAlreadyExists{Login: login}
}

func (e *ErrorUserAlreadyExists) Error() string {
	return fmt.Sprintf("User with login %s already exists", e.Login)
}
