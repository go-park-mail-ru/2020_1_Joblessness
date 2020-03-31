// Code generated by MockGen. DO NOT EDIT.
// Source: joblessness/haha/auth/interfaces (interfaces: AuthRepository)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	models "joblessness/haha/models"
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

// ChangeOrganization mocks base method
func (m *MockAuthRepository) ChangeOrganization(arg0 models.Organization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeOrganization", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeOrganization indicates an expected call of ChangeOrganization
func (mr *MockAuthRepositoryMockRecorder) ChangeOrganization(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeOrganization", reflect.TypeOf((*MockAuthRepository)(nil).ChangeOrganization), arg0)
}

// ChangePerson mocks base method
func (m *MockAuthRepository) ChangePerson(arg0 models.Person) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePerson", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePerson indicates an expected call of ChangePerson
func (mr *MockAuthRepositoryMockRecorder) ChangePerson(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePerson", reflect.TypeOf((*MockAuthRepository)(nil).ChangePerson), arg0)
}

// CreateOrganization mocks base method
func (m *MockAuthRepository) CreateOrganization(arg0 *models.Organization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrganization", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrganization indicates an expected call of CreateOrganization
func (mr *MockAuthRepositoryMockRecorder) CreateOrganization(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrganization", reflect.TypeOf((*MockAuthRepository)(nil).CreateOrganization), arg0)
}

// CreatePerson mocks base method
func (m *MockAuthRepository) CreatePerson(arg0 *models.Person) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePerson", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePerson indicates an expected call of CreatePerson
func (mr *MockAuthRepositoryMockRecorder) CreatePerson(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePerson", reflect.TypeOf((*MockAuthRepository)(nil).CreatePerson), arg0)
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

// GetListOfOrgs mocks base method
func (m *MockAuthRepository) GetListOfOrgs(arg0 int) ([]models.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListOfOrgs", arg0)
	ret0, _ := ret[0].([]models.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListOfOrgs indicates an expected call of GetListOfOrgs
func (mr *MockAuthRepositoryMockRecorder) GetListOfOrgs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListOfOrgs", reflect.TypeOf((*MockAuthRepository)(nil).GetListOfOrgs), arg0)
}

// GetOrganization mocks base method
func (m *MockAuthRepository) GetOrganization(arg0 uint64) (*models.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganization", arg0)
	ret0, _ := ret[0].(*models.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganization indicates an expected call of GetOrganization
func (mr *MockAuthRepositoryMockRecorder) GetOrganization(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganization", reflect.TypeOf((*MockAuthRepository)(nil).GetOrganization), arg0)
}

// GetPerson mocks base method
func (m *MockAuthRepository) GetPerson(arg0 uint64) (*models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPerson", arg0)
	ret0, _ := ret[0].(*models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPerson indicates an expected call of GetPerson
func (mr *MockAuthRepositoryMockRecorder) GetPerson(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPerson", reflect.TypeOf((*MockAuthRepository)(nil).GetPerson), arg0)
}

// GetUserFavorite mocks base method
func (m *MockAuthRepository) GetUserFavorite(arg0 uint64) (models.Favorites, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFavorite", arg0)
	ret0, _ := ret[0].(models.Favorites)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFavorite indicates an expected call of GetUserFavorite
func (mr *MockAuthRepositoryMockRecorder) GetUserFavorite(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFavorite", reflect.TypeOf((*MockAuthRepository)(nil).GetUserFavorite), arg0)
}

// IsOrganization mocks base method
func (m *MockAuthRepository) IsOrganization(arg0 uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOrganization", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsOrganization indicates an expected call of IsOrganization
func (mr *MockAuthRepositoryMockRecorder) IsOrganization(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOrganization", reflect.TypeOf((*MockAuthRepository)(nil).IsOrganization), arg0)
}

// IsPerson mocks base method
func (m *MockAuthRepository) IsPerson(arg0 uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPerson", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPerson indicates an expected call of IsPerson
func (mr *MockAuthRepositoryMockRecorder) IsPerson(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPerson", reflect.TypeOf((*MockAuthRepository)(nil).IsPerson), arg0)
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

// SaveAvatarLink mocks base method
func (m *MockAuthRepository) SaveAvatarLink(arg0 string, arg1 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAvatarLink", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAvatarLink indicates an expected call of SaveAvatarLink
func (mr *MockAuthRepositoryMockRecorder) SaveAvatarLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAvatarLink", reflect.TypeOf((*MockAuthRepository)(nil).SaveAvatarLink), arg0, arg1)
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

// SetOrDeleteLike mocks base method
func (m *MockAuthRepository) SetOrDeleteLike(arg0, arg1 uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOrDeleteLike", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetOrDeleteLike indicates an expected call of SetOrDeleteLike
func (mr *MockAuthRepositoryMockRecorder) SetOrDeleteLike(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOrDeleteLike", reflect.TypeOf((*MockAuthRepository)(nil).SetOrDeleteLike), arg0, arg1)
}
