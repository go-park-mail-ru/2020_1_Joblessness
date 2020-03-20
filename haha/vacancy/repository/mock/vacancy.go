// Code generated by MockGen. DO NOT EDIT.
// Source: joblessness/haha/vacancy (interfaces: VacancyRepository)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	models "joblessness/haha/models"
	reflect "reflect"
)

// MockVacancyRepository is a mock of VacancyRepository interface
type MockVacancyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVacancyRepositoryMockRecorder
}

// MockVacancyRepositoryMockRecorder is the mock recorder for MockVacancyRepository
type MockVacancyRepositoryMockRecorder struct {
	mock *MockVacancyRepository
}

// NewMockVacancyRepository creates a new mock instance
func NewMockVacancyRepository(ctrl *gomock.Controller) *MockVacancyRepository {
	mock := &MockVacancyRepository{ctrl: ctrl}
	mock.recorder = &MockVacancyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVacancyRepository) EXPECT() *MockVacancyRepositoryMockRecorder {
	return m.recorder
}

// ChangeVacancy mocks base method
func (m *MockVacancyRepository) ChangeVacancy(arg0 *models.Vacancy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeVacancy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeVacancy indicates an expected call of ChangeVacancy
func (mr *MockVacancyRepositoryMockRecorder) ChangeVacancy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeVacancy", reflect.TypeOf((*MockVacancyRepository)(nil).ChangeVacancy), arg0)
}

// CreateVacancy mocks base method
func (m *MockVacancyRepository) CreateVacancy(arg0 *models.Vacancy) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVacancy", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVacancy indicates an expected call of CreateVacancy
func (mr *MockVacancyRepositoryMockRecorder) CreateVacancy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVacancy", reflect.TypeOf((*MockVacancyRepository)(nil).CreateVacancy), arg0)
}

// DeleteVacancy mocks base method
func (m *MockVacancyRepository) DeleteVacancy(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVacancy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVacancy indicates an expected call of DeleteVacancy
func (mr *MockVacancyRepositoryMockRecorder) DeleteVacancy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVacancy", reflect.TypeOf((*MockVacancyRepository)(nil).DeleteVacancy), arg0)
}

// GetVacancies mocks base method
func (m *MockVacancyRepository) GetVacancies() ([]models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacancies")
	ret0, _ := ret[0].([]models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacancies indicates an expected call of GetVacancies
func (mr *MockVacancyRepositoryMockRecorder) GetVacancies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacancies", reflect.TypeOf((*MockVacancyRepository)(nil).GetVacancies))
}

// GetVacancy mocks base method
func (m *MockVacancyRepository) GetVacancy(arg0 uint64) (*models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacancy", arg0)
	ret0, _ := ret[0].(*models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacancy indicates an expected call of GetVacancy
func (mr *MockVacancyRepositoryMockRecorder) GetVacancy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacancy", reflect.TypeOf((*MockVacancyRepository)(nil).GetVacancy), arg0)
}
