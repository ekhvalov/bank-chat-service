// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/usecase_mock.gen.go -typed -package=getchatsmocks
//
// Package getchatsmocks is a generated GoMock package.
package getchatsmocks

import (
	context "context"
	reflect "reflect"

	chatsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/chats"
	types "github.com/ekhvalov/bank-chat-service/internal/types"
	gomock "go.uber.org/mock/gomock"
)

// MockchatsRepository is a mock of chatsRepository interface.
type MockchatsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockchatsRepositoryMockRecorder
}

// MockchatsRepositoryMockRecorder is the mock recorder for MockchatsRepository.
type MockchatsRepositoryMockRecorder struct {
	mock *MockchatsRepository
}

// NewMockchatsRepository creates a new mock instance.
func NewMockchatsRepository(ctrl *gomock.Controller) *MockchatsRepository {
	mock := &MockchatsRepository{ctrl: ctrl}
	mock.recorder = &MockchatsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockchatsRepository) EXPECT() *MockchatsRepositoryMockRecorder {
	return m.recorder
}

// GetOpenProblemChatsForManager mocks base method.
func (m *MockchatsRepository) GetOpenProblemChatsForManager(ctx context.Context, managerID types.UserID) ([]chatsrepo.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenProblemChatsForManager", ctx, managerID)
	ret0, _ := ret[0].([]chatsrepo.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenProblemChatsForManager indicates an expected call of GetOpenProblemChatsForManager.
func (mr *MockchatsRepositoryMockRecorder) GetOpenProblemChatsForManager(ctx, managerID any) *chatsRepositoryGetOpenProblemChatsForManagerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenProblemChatsForManager", reflect.TypeOf((*MockchatsRepository)(nil).GetOpenProblemChatsForManager), ctx, managerID)
	return &chatsRepositoryGetOpenProblemChatsForManagerCall{Call: call}
}

// chatsRepositoryGetOpenProblemChatsForManagerCall wrap *gomock.Call
type chatsRepositoryGetOpenProblemChatsForManagerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *chatsRepositoryGetOpenProblemChatsForManagerCall) Return(arg0 []chatsrepo.Chat, arg1 error) *chatsRepositoryGetOpenProblemChatsForManagerCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *chatsRepositoryGetOpenProblemChatsForManagerCall) Do(f func(context.Context, types.UserID) ([]chatsrepo.Chat, error)) *chatsRepositoryGetOpenProblemChatsForManagerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *chatsRepositoryGetOpenProblemChatsForManagerCall) DoAndReturn(f func(context.Context, types.UserID) ([]chatsrepo.Chat, error)) *chatsRepositoryGetOpenProblemChatsForManagerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
