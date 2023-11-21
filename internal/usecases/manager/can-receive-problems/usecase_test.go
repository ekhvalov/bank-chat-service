package canreceiveproblems_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/ekhvalov/bank-chat-service/internal/testingh"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	canreceiveproblemsmocks "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *canreceiveproblemsmocks.MockmanagerLoadService
	mPoolMock *canreceiveproblemsmocks.MockmanagerPool
	uCase     canreceiveproblems.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mPoolMock = canreceiveproblemsmocks.NewMockmanagerPool(s.ctrl)
	s.mLoadMock = canreceiveproblemsmocks.NewMockmanagerLoadService(s.ctrl)

	var err error
	s.uCase, err = canreceiveproblems.New(canreceiveproblems.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestUseCase_CreateError() {
	useCase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(nil, nil))

	s.Require().Error(err)
	s.Assert().Empty(useCase)
}

func (s *UseCaseSuite) TestUseCase_InvalidRequest() {
	// Arrange
	request := canreceiveproblems.Request{ID: types.RequestID{}, ManagerID: types.UserID{}}

	// Action
	result, err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().ErrorIs(err, canreceiveproblems.ErrInvalidRequest)
	s.Assert().False(result.Result)
}

func (s *UseCaseSuite) TestUseCase_ManagerPoolError() {
	// Arrange
	request := canreceiveproblems.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	errPool := errors.New("pool error")
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, errPool)

	// Action
	result, err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().ErrorIs(err, canreceiveproblems.ErrManagerPool)
	s.Assert().False(result.Result)
}

func (s *UseCaseSuite) TestUseCase_ManagerAlreadyInPool() {
	// Arrange
	request := canreceiveproblems.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(true, nil)

	// Action
	result, err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().NoError(err)
	s.Assert().False(result.Result)
}

func (s *UseCaseSuite) TestUseCase_ManagerLoadServiceError() {
	// Arrange
	request := canreceiveproblems.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, nil)
	errLoadService := errors.New("load service error")
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, request.ManagerID).Return(false, errLoadService)

	// Action
	result, err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().ErrorIs(err, canreceiveproblems.ErrManagerLoadService)
	s.Assert().False(result.Result)
}

func (s *UseCaseSuite) TestUseCase_ManagerCanTakeProblem() {
	// Arrange
	request := canreceiveproblems.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, request.ManagerID).Return(true, nil)

	// Action
	result, err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().NoError(err)
	s.Assert().True(result.Result)
}

func (s *UseCaseSuite) TestUseCase_ManagerCanNotTakeProblem() {
	// Arrange
	request := canreceiveproblems.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID()}
	s.mPoolMock.EXPECT().Contains(s.Ctx, request.ManagerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, request.ManagerID).Return(false, nil)

	// Action
	result, err := s.uCase.Handle(s.Ctx, request)

	// Assert
	s.Require().NoError(err)
	s.Assert().False(result.Result)
}
