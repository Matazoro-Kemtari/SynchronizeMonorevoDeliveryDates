// Code generated by MockGen. DO NOT EDIT.
// Source: .\proposition_fetch_case.go

// Package mock_proposition_fetch_case is a generated GoMock package.
package mock_proposition_fetch_case

import (
	proposition_fetch_case "SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFetchingExecutor is a mock of FetchingExecutor interface.
type MockFetchingExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockFetchingExecutorMockRecorder
}

// MockFetchingExecutorMockRecorder is the mock recorder for MockFetchingExecutor.
type MockFetchingExecutorMockRecorder struct {
	mock *MockFetchingExecutor
}

// NewMockFetchingExecutor creates a new mock instance.
func NewMockFetchingExecutor(ctrl *gomock.Controller) *MockFetchingExecutor {
	mock := &MockFetchingExecutor{ctrl: ctrl}
	mock.recorder = &MockFetchingExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFetchingExecutor) EXPECT() *MockFetchingExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockFetchingExecutor) Execute() ([]proposition_fetch_case.FetchedPropositionDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute")
	ret0, _ := ret[0].([]proposition_fetch_case.FetchedPropositionDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockFetchingExecutorMockRecorder) Execute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockFetchingExecutor)(nil).Execute))
}
