// Code generated by MockGen. DO NOT EDIT.
// Source: joblessness/haha/auth/interfaces (interfaces: AuthRepository)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAuthRepository is a mock of AuthRepository interface
type MockAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepositoryMockRecorder
}

// MockAuthRepositoryMockRecorder is the mock recorder for MockAuthRepository
type MockAuthRepositoryMockRecorder struct {
	mock *MockAuthRepository
}

// NewMockAuthRepository creates a new mock instance
func NewMockAuthRepository(ctrl *gomock.Controller) *MockAuthRepository {
	mock := &MockAuthRepository{ctrl: ctrl}
	mock.recorder = &MockAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthRepository) EXPECT() *MockAuthRepositoryMockRecorder {
	return m.recorder
}

// DoesUserExists mocks base method
func (m *MockAuthRepository) DoesUserExists(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoesUserExists", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DoesUserExists indicates an expected call of DoesUserExists
func (mr *MockAuthRepositoryMockRecorder) DoesUserExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoesUserExists", reflect.TypeOf((*MockAuthRepository)(nil).DoesUserExists), arg0)
}

// GetRole mocks base method
func (m *MockAuthRepository) GetRole(arg0 uint64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole
func (mr *MockAuthRepositoryMockRecorder) GetRole(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockAuthRepository)(nil).GetRole), arg0)
}

// Login mocks base method
func (m *MockAuthRepository) Login(arg0, arg1, arg2 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login
func (mr *MockAuthRepositoryMockRecorder) Login(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthRepository)(nil).Login), arg0, arg1, arg2)
}

// Logout mocks base method
func (m *MockAuthRepository) Logout(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout
func (mr *MockAuthRepositoryMockRecorder) Logout(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthRepository)(nil).Logout), arg0)
}

// RegisterOrganization mocks base method
func (m *MockAuthRepository) RegisterOrganization(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterOrganization", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterOrganization indicates an expected call of RegisterOrganization
func (mr *MockAuthRepositoryMockRecorder) RegisterOrganization(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterOrganization", reflect.TypeOf((*MockAuthRepository)(nil).RegisterOrganization), arg0, arg1, arg2)
}

// RegisterPerson mocks base method
func (m *MockAuthRepository) RegisterPerson(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterPerson", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterPerson indicates an expected call of RegisterPerson
func (mr *MockAuthRepositoryMockRecorder) RegisterPerson(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPerson", reflect.TypeOf((*MockAuthRepository)(nil).RegisterPerson), arg0, arg1, arg2)
}

// SessionExists mocks base method
func (m *MockAuthRepository) SessionExists(arg0 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SessionExists", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SessionExists indicates an expected call of SessionExists
func (mr *MockAuthRepositoryMockRecorder) SessionExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SessionExists", reflect.TypeOf((*MockAuthRepository)(nil).SessionExists), arg0)
}
