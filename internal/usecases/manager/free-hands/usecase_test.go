package freehands_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/ekhvalov/bank-chat-service/internal/testingh"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
	freehandsmocks "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *freehandsmocks.MockmanagerLoadService
	mPoolMock *freehandsmocks.MockmanagerPool
	uCase     freehands.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mPoolMock = freehandsmocks.NewMockmanagerPool(s.ctrl)
	s.mLoadMock = freehandsmocks.NewMockmanagerLoadService(s.ctrl)

	var err error
	s.uCase, err = freehands.New(freehands.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestUseCase_CreateError() {
	useCase, err := freehands.New(freehands.NewOptions(nil, nil))

	s.Require().Error(err)
	s.Assert().Empty(useCase)
}

func (s *UseCaseSuite) TestUseCase_InvalidRequest() {
	// Arrange
	request := freehands.Request{ID: types.RequestID{}, ManagerID: types.UserID{}}

	// Action
	err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().ErrorIs(err, freehands.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestUseCase_ManagerPool_ContainsError() {
	// Arrange
	request := freehands.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	errPool := errors.New("pool error")
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, errPool)

	// Action
	err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestUseCase_ManagerAlreadyInPool() {
	// Arrange
	request := freehands.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(true, nil)

	// Action
	err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().NoError(err)
}

func (s *UseCaseSuite) TestUseCase_ManagerLoadServiceError() {
	// Arrange
	request := freehands.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, nil)
	errLoadService := errors.New("load service error")
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, request.ManagerID).Return(false, errLoadService)

	// Action
	err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestUseCase_ManagerPool_PutError() {
	// Arrange
	request := freehands.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, request.ManagerID).Return(true, nil)
	errPool := errors.New("pool error")
	s.mPoolMock.EXPECT().Put(s.Ctx, request.ManagerID).Return(errPool)

	// Action
	err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestUseCase_ManagerCanTakeProblem() {
	// Arrange
	request := freehands.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, request.ManagerID).Return(true, nil)
	s.mPoolMock.EXPECT().Put(s.Ctx, request.ManagerID).Return(nil)

	// Action
	err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().NoError(err)
}

func (s *UseCaseSuite) TestUseCase_ManagerCanNotTakeProblem() {
	// Arrange
	request := freehands.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, request.ManagerID).Return(false, nil)

	// Action
	err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().ErrorIs(err, freehands.ErrManagerOverloaded)
}
