package managerload_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	managerload "github.com/ekhvalov/bank-chat-service/internal/services/manager-load"
	managerloadmocks "github.com/ekhvalov/bank-chat-service/internal/services/manager-load/mocks"
	"github.com/ekhvalov/bank-chat-service/internal/testingh"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const maxProblemsAtTime = 30

type ServiceSuite struct {
	testingh.ContextSuite

	ctrl *gomock.Controller

	problemsRepo *managerloadmocks.MockproblemsRepository
	managerLoad  *managerload.Service
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.problemsRepo = managerloadmocks.NewMockproblemsRepository(s.ctrl)

	var err error
	s.managerLoad, err = managerload.New(managerload.NewOptions(maxProblemsAtTime, s.problemsRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *ServiceSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *ServiceSuite) TestInvalidOptions() {
	s.Run("problems limit too low", func() {
		// Action
		_, err := managerload.New(managerload.NewOptions(0, s.problemsRepo))

		// Assert
		s.Require().Error(err)
	})

	s.Run("problems limit too high", func() {
		// Action
		_, err := managerload.New(managerload.NewOptions(maxProblemsAtTime+1, s.problemsRepo))

		// Assert
		s.Require().Error(err)
	})

	s.Run("nil problems repo", func() {
		// Action
		_, err := managerload.New(managerload.NewOptions(maxProblemsAtTime, nil))

		// Assert
		s.Require().Error(err)
	})
}

func (s *ServiceSuite) TestService_CanManagerTakeProblem_RepoError() {
	// Arrange
	managerID := types.NewUserID()
	errRepo := errors.New("repo error")
	s.problemsRepo.EXPECT().CountManagerOpenProblems(s.Ctx, managerID).Return(0, errRepo)

	// Action
	ok, err := s.managerLoad.CanManagerTakeProblem(s.Ctx, managerID)

	// Assert
	s.Require().Error(err)
	s.Assert().False(ok)
}

func (s *ServiceSuite) TestService_CanManagerTakeProblem() {
	tests := []struct {
		openProblems   int
		expectedResult bool
	}{
		{openProblems: 0, expectedResult: true},
		{openProblems: 1, expectedResult: true},
		{openProblems: maxProblemsAtTime - 1, expectedResult: true},
		{openProblems: maxProblemsAtTime, expectedResult: false},
		{openProblems: maxProblemsAtTime + 1, expectedResult: false},
	}
	for _, tt := range tests {
		s.Run(fmt.Sprintf("open problems %d", tt.openProblems), func() {
			// Arrange
			managerID := types.NewUserID()
			s.problemsRepo.EXPECT().CountManagerOpenProblems(s.Ctx, managerID).Return(tt.openProblems, nil)

			// Action
			result, err := s.managerLoad.CanManagerTakeProblem(s.Ctx, managerID)

			// Assert
			s.Require().NoError(err)
			s.Assert().Equal(tt.expectedResult, result)
		})
	}
}
