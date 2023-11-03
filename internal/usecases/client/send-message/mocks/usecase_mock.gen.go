// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks
//
// Package sendmessagemocks is a generated GoMock package.
package sendmessagemocks

import (
	context "context"
	reflect "reflect"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
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

// CreateIfNotExists mocks base method.
func (m *MockchatsRepository) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIfNotExists", ctx, userID)
	ret0, _ := ret[0].(types.ChatID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIfNotExists indicates an expected call of CreateIfNotExists.
func (mr *MockchatsRepositoryMockRecorder) CreateIfNotExists(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIfNotExists", reflect.TypeOf((*MockchatsRepository)(nil).CreateIfNotExists), ctx, userID)
}

// MockmessagesRepository is a mock of messagesRepository interface.
type MockmessagesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockmessagesRepositoryMockRecorder
}

// MockmessagesRepositoryMockRecorder is the mock recorder for MockmessagesRepository.
type MockmessagesRepositoryMockRecorder struct {
	mock *MockmessagesRepository
}

// NewMockmessagesRepository creates a new mock instance.
func NewMockmessagesRepository(ctrl *gomock.Controller) *MockmessagesRepository {
	mock := &MockmessagesRepository{ctrl: ctrl}
	mock.recorder = &MockmessagesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmessagesRepository) EXPECT() *MockmessagesRepositoryMockRecorder {
	return m.recorder
}

// CreateClientVisible mocks base method.
func (m *MockmessagesRepository) CreateClientVisible(ctx context.Context, reqID types.RequestID, problemID types.ProblemID, chatID types.ChatID, authorID types.UserID, msgBody string) (*messagesrepo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClientVisible", ctx, reqID, problemID, chatID, authorID, msgBody)
	ret0, _ := ret[0].(*messagesrepo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateClientVisible indicates an expected call of CreateClientVisible.
func (mr *MockmessagesRepositoryMockRecorder) CreateClientVisible(ctx, reqID, problemID, chatID, authorID, msgBody any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClientVisible", reflect.TypeOf((*MockmessagesRepository)(nil).CreateClientVisible), ctx, reqID, problemID, chatID, authorID, msgBody)
}

// GetMessageByRequestID mocks base method.
func (m *MockmessagesRepository) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*messagesrepo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageByRequestID", ctx, reqID)
	ret0, _ := ret[0].(*messagesrepo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageByRequestID indicates an expected call of GetMessageByRequestID.
func (mr *MockmessagesRepositoryMockRecorder) GetMessageByRequestID(ctx, reqID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageByRequestID", reflect.TypeOf((*MockmessagesRepository)(nil).GetMessageByRequestID), ctx, reqID)
}

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

// CreateIfNotExists mocks base method.
func (m *MockproblemsRepository) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIfNotExists", ctx, chatID)
	ret0, _ := ret[0].(types.ProblemID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIfNotExists indicates an expected call of CreateIfNotExists.
func (mr *MockproblemsRepositoryMockRecorder) CreateIfNotExists(ctx, chatID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIfNotExists", reflect.TypeOf((*MockproblemsRepository)(nil).CreateIfNotExists), ctx, chatID)
}

// Mocktransactor is a mock of transactor interface.
type Mocktransactor struct {
	ctrl     *gomock.Controller
	recorder *MocktransactorMockRecorder
}

// MocktransactorMockRecorder is the mock recorder for Mocktransactor.
type MocktransactorMockRecorder struct {
	mock *Mocktransactor
}

// NewMocktransactor creates a new mock instance.
func NewMocktransactor(ctrl *gomock.Controller) *Mocktransactor {
	mock := &Mocktransactor{ctrl: ctrl}
	mock.recorder = &MocktransactorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktransactor) EXPECT() *MocktransactorMockRecorder {
	return m.recorder
}

// RunInTx mocks base method.
func (m *Mocktransactor) RunInTx(ctx context.Context, f func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTx", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTx indicates an expected call of RunInTx.
func (mr *MocktransactorMockRecorder) RunInTx(ctx, f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTx", reflect.TypeOf((*Mocktransactor)(nil).RunInTx), ctx, f)
}
