// Code generated by MockGen. DO NOT EDIT.
// Source: joblessness/haha/summary/interfaces (interfaces: SummaryUseCase)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	models "joblessness/haha/models"
	reflect "reflect"
)

// MockSummaryUseCase is a mock of SummaryUseCase interface
type MockSummaryUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSummaryUseCaseMockRecorder
}

// MockSummaryUseCaseMockRecorder is the mock recorder for MockSummaryUseCase
type MockSummaryUseCaseMockRecorder struct {
	mock *MockSummaryUseCase
}

// NewMockSummaryUseCase creates a new mock instance
func NewMockSummaryUseCase(ctrl *gomock.Controller) *MockSummaryUseCase {
	mock := &MockSummaryUseCase{ctrl: ctrl}
	mock.recorder = &MockSummaryUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSummaryUseCase) EXPECT() *MockSummaryUseCaseMockRecorder {
	return m.recorder
}

// ChangeSummary mocks base method
func (m *MockSummaryUseCase) ChangeSummary(arg0 *models.Summary) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeSummary", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeSummary indicates an expected call of ChangeSummary
func (mr *MockSummaryUseCaseMockRecorder) ChangeSummary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeSummary", reflect.TypeOf((*MockSummaryUseCase)(nil).ChangeSummary), arg0)
}

// CreateSummary mocks base method
func (m *MockSummaryUseCase) CreateSummary(arg0 *models.Summary) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSummary", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSummary indicates an expected call of CreateSummary
func (mr *MockSummaryUseCaseMockRecorder) CreateSummary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSummary", reflect.TypeOf((*MockSummaryUseCase)(nil).CreateSummary), arg0)
}

// DeleteSummary mocks base method
func (m *MockSummaryUseCase) DeleteSummary(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSummary", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSummary indicates an expected call of DeleteSummary
func (mr *MockSummaryUseCaseMockRecorder) DeleteSummary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSummary", reflect.TypeOf((*MockSummaryUseCase)(nil).DeleteSummary), arg0)
}

// GetAllSummaries mocks base method
func (m *MockSummaryUseCase) GetAllSummaries() ([]models.Summary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSummaries")
	ret0, _ := ret[0].([]models.Summary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSummaries indicates an expected call of GetAllSummaries
func (mr *MockSummaryUseCaseMockRecorder) GetAllSummaries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSummaries", reflect.TypeOf((*MockSummaryUseCase)(nil).GetAllSummaries))
}

// GetOrgSummaries mocks base method
func (m *MockSummaryUseCase) GetOrgSummaries(arg0 uint64) (models.OrgSummaries, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrgSummaries", arg0)
	ret0, _ := ret[0].(models.OrgSummaries)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrgSummaries indicates an expected call of GetOrgSummaries
func (mr *MockSummaryUseCaseMockRecorder) GetOrgSummaries(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrgSummaries", reflect.TypeOf((*MockSummaryUseCase)(nil).GetOrgSummaries), arg0)
}

// GetSummary mocks base method
func (m *MockSummaryUseCase) GetSummary(arg0 uint64) (*models.Summary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSummary", arg0)
	ret0, _ := ret[0].(*models.Summary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSummary indicates an expected call of GetSummary
func (mr *MockSummaryUseCaseMockRecorder) GetSummary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSummary", reflect.TypeOf((*MockSummaryUseCase)(nil).GetSummary), arg0)
}

// GetUserSummaries mocks base method
func (m *MockSummaryUseCase) GetUserSummaries(arg0 uint64) ([]models.Summary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSummaries", arg0)
	ret0, _ := ret[0].([]models.Summary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSummaries indicates an expected call of GetUserSummaries
func (mr *MockSummaryUseCaseMockRecorder) GetUserSummaries(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSummaries", reflect.TypeOf((*MockSummaryUseCase)(nil).GetUserSummaries), arg0)
}

// ResponseSummary mocks base method
func (m *MockSummaryUseCase) ResponseSummary(arg0 *models.SendSummary) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResponseSummary", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResponseSummary indicates an expected call of ResponseSummary
func (mr *MockSummaryUseCaseMockRecorder) ResponseSummary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseSummary", reflect.TypeOf((*MockSummaryUseCase)(nil).ResponseSummary), arg0)
}

// SendSummary mocks base method
func (m *MockSummaryUseCase) SendSummary(arg0 *models.SendSummary) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSummary", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSummary indicates an expected call of SendSummary
func (mr *MockSummaryUseCaseMockRecorder) SendSummary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSummary", reflect.TypeOf((*MockSummaryUseCase)(nil).SendSummary), arg0)
}
