// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go
//
// Generated by this command:
//
//	mockgen -source=handlers.go -destination=mocks/handlers_mocks.gen.go -package=managerv1mocks
//
// Package managerv1mocks is a generated GoMock package.
package managerv1mocks

import (
	context "context"
	reflect "reflect"

	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
	gomock "go.uber.org/mock/gomock"
)

// MockcanReceiveProblemsUsecase is a mock of canReceiveProblemsUsecase interface.
type MockcanReceiveProblemsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockcanReceiveProblemsUsecaseMockRecorder
}

// MockcanReceiveProblemsUsecaseMockRecorder is the mock recorder for MockcanReceiveProblemsUsecase.
type MockcanReceiveProblemsUsecaseMockRecorder struct {
	mock *MockcanReceiveProblemsUsecase
}

// NewMockcanReceiveProblemsUsecase creates a new mock instance.
func NewMockcanReceiveProblemsUsecase(ctrl *gomock.Controller) *MockcanReceiveProblemsUsecase {
	mock := &MockcanReceiveProblemsUsecase{ctrl: ctrl}
	mock.recorder = &MockcanReceiveProblemsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcanReceiveProblemsUsecase) EXPECT() *MockcanReceiveProblemsUsecaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockcanReceiveProblemsUsecase) Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(canreceiveproblems.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockcanReceiveProblemsUsecaseMockRecorder) Handle(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockcanReceiveProblemsUsecase)(nil).Handle), ctx, req)
}

// MockfreeHandsUsecase is a mock of freeHandsUsecase interface.
type MockfreeHandsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockfreeHandsUsecaseMockRecorder
}

// MockfreeHandsUsecaseMockRecorder is the mock recorder for MockfreeHandsUsecase.
type MockfreeHandsUsecaseMockRecorder struct {
	mock *MockfreeHandsUsecase
}

// NewMockfreeHandsUsecase creates a new mock instance.
func NewMockfreeHandsUsecase(ctrl *gomock.Controller) *MockfreeHandsUsecase {
	mock := &MockfreeHandsUsecase{ctrl: ctrl}
	mock.recorder = &MockfreeHandsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockfreeHandsUsecase) EXPECT() *MockfreeHandsUsecaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockfreeHandsUsecase) Handle(ctx context.Context, req freehands.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockfreeHandsUsecaseMockRecorder) Handle(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockfreeHandsUsecase)(nil).Handle), ctx, req)
}
