// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mocks/service_mock.gen.go -typed -package=afcverdictsprocessormocks
//
// Package afcverdictsprocessormocks is a generated GoMock package.
package afcverdictsprocessormocks

import (
	context "context"
	reflect "reflect"
	time "time"

	types "github.com/ekhvalov/bank-chat-service/internal/types"
	gomock "go.uber.org/mock/gomock"
)

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

// BlockMessage mocks base method.
func (m *MockmessagesRepository) BlockMessage(ctx context.Context, msgID types.MessageID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockMessage", ctx, msgID)
	ret0, _ := ret[0].(error)
	return ret0
}

// BlockMessage indicates an expected call of BlockMessage.
func (mr *MockmessagesRepositoryMockRecorder) BlockMessage(ctx, msgID any) *messagesRepositoryBlockMessageCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockMessage", reflect.TypeOf((*MockmessagesRepository)(nil).BlockMessage), ctx, msgID)
	return &messagesRepositoryBlockMessageCall{Call: call}
}

// messagesRepositoryBlockMessageCall wrap *gomock.Call
type messagesRepositoryBlockMessageCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *messagesRepositoryBlockMessageCall) Return(arg0 error) *messagesRepositoryBlockMessageCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *messagesRepositoryBlockMessageCall) Do(f func(context.Context, types.MessageID) error) *messagesRepositoryBlockMessageCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *messagesRepositoryBlockMessageCall) DoAndReturn(f func(context.Context, types.MessageID) error) *messagesRepositoryBlockMessageCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MarkAsVisibleForManager mocks base method.
func (m *MockmessagesRepository) MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsVisibleForManager", ctx, msgID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkAsVisibleForManager indicates an expected call of MarkAsVisibleForManager.
func (mr *MockmessagesRepositoryMockRecorder) MarkAsVisibleForManager(ctx, msgID any) *messagesRepositoryMarkAsVisibleForManagerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsVisibleForManager", reflect.TypeOf((*MockmessagesRepository)(nil).MarkAsVisibleForManager), ctx, msgID)
	return &messagesRepositoryMarkAsVisibleForManagerCall{Call: call}
}

// messagesRepositoryMarkAsVisibleForManagerCall wrap *gomock.Call
type messagesRepositoryMarkAsVisibleForManagerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *messagesRepositoryMarkAsVisibleForManagerCall) Return(arg0 error) *messagesRepositoryMarkAsVisibleForManagerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *messagesRepositoryMarkAsVisibleForManagerCall) Do(f func(context.Context, types.MessageID) error) *messagesRepositoryMarkAsVisibleForManagerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *messagesRepositoryMarkAsVisibleForManagerCall) DoAndReturn(f func(context.Context, types.MessageID) error) *messagesRepositoryMarkAsVisibleForManagerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockoutboxService is a mock of outboxService interface.
type MockoutboxService struct {
	ctrl     *gomock.Controller
	recorder *MockoutboxServiceMockRecorder
}

// MockoutboxServiceMockRecorder is the mock recorder for MockoutboxService.
type MockoutboxServiceMockRecorder struct {
	mock *MockoutboxService
}

// NewMockoutboxService creates a new mock instance.
func NewMockoutboxService(ctrl *gomock.Controller) *MockoutboxService {
	mock := &MockoutboxService{ctrl: ctrl}
	mock.recorder = &MockoutboxServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockoutboxService) EXPECT() *MockoutboxServiceMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockoutboxService) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, name, payload, availableAt)
	ret0, _ := ret[0].(types.JobID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockoutboxServiceMockRecorder) Put(ctx, name, payload, availableAt any) *outboxServicePutCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockoutboxService)(nil).Put), ctx, name, payload, availableAt)
	return &outboxServicePutCall{Call: call}
}

// outboxServicePutCall wrap *gomock.Call
type outboxServicePutCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *outboxServicePutCall) Return(arg0 types.JobID, arg1 error) *outboxServicePutCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *outboxServicePutCall) Do(f func(context.Context, string, string, time.Time) (types.JobID, error)) *outboxServicePutCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *outboxServicePutCall) DoAndReturn(f func(context.Context, string, string, time.Time) (types.JobID, error)) *outboxServicePutCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MocktransactorMockRecorder) RunInTx(ctx, f any) *transactorRunInTxCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTx", reflect.TypeOf((*Mocktransactor)(nil).RunInTx), ctx, f)
	return &transactorRunInTxCall{Call: call}
}

// transactorRunInTxCall wrap *gomock.Call
type transactorRunInTxCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *transactorRunInTxCall) Return(arg0 error) *transactorRunInTxCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *transactorRunInTxCall) Do(f func(context.Context, func(context.Context) error) error) *transactorRunInTxCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *transactorRunInTxCall) DoAndReturn(f func(context.Context, func(context.Context) error) error) *transactorRunInTxCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
