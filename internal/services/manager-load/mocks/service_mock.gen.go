// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mocks/service_mock.gen.go -package=managerloadmocks
//
// Package managerloadmocks is a generated GoMock package.
package managerloadmocks

import (
	context "context"
	reflect "reflect"

	types "github.com/ekhvalov/bank-chat-service/internal/types"
	gomock "go.uber.org/mock/gomock"
)

// MockproblemsRepository is a mock of problemsRepository interface.
type MockproblemsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockproblemsRepositoryMockRecorder
}

// MockproblemsRepositoryMockRecorder is the mock recorder for MockproblemsRepository.
type MockproblemsRepositoryMockRecorder struct {
	mock *MockproblemsRepository
}

// NewMockproblemsRepository creates a new mock instance.
func NewMockproblemsRepository(ctrl *gomock.Controller) *MockproblemsRepository {
	mock := &MockproblemsRepository{ctrl: ctrl}
	mock.recorder = &MockproblemsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproblemsRepository) EXPECT() *MockproblemsRepositoryMockRecorder {
	return m.recorder
}

// CountManagerOpenProblems mocks base method.
func (m *MockproblemsRepository) CountManagerOpenProblems(ctx context.Context, managerID types.UserID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountManagerOpenProblems", ctx, managerID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountManagerOpenProblems indicates an expected call of CountManagerOpenProblems.
func (mr *MockproblemsRepositoryMockRecorder) CountManagerOpenProblems(ctx, managerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountManagerOpenProblems", reflect.TypeOf((*MockproblemsRepository)(nil).CountManagerOpenProblems), ctx, managerID)
}