package getchathistory_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/ekhvalov/bank-chat-service/internal/cursor"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	"github.com/ekhvalov/bank-chat-service/internal/testingh"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	getchathistory "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chat-history"
	getchathistorymocks "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chat-history/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl         *gomock.Controller
	messagesRepo *getchathistorymocks.MockmessagesRepository
	problemsRepo *getchathistorymocks.MockproblemsRepo
	uCase        getchathistory.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.messagesRepo = getchathistorymocks.NewMockmessagesRepository(s.ctrl)
	s.problemsRepo = getchathistorymocks.NewMockproblemsRepo(s.ctrl)

	var err error
	s.uCase, err = getchathistory.New(getchathistory.NewOptions(s.messagesRepo, s.problemsRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestHandle_RequestValidationError() {
	// Arrange.
	req := getchathistory.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidRequest)
	s.Empty(resp.Messages)
}

func (s *UseCaseSuite) TestHandle_CursorDecodingError() {
	// Arrange.
	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
		ChatID:    types.NewChatID(),
		Cursor:    "eyJwYWdlX3NpemUiOjEwMA==", // {"page_size":100
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidCursor)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestHandle_ProblemsRepoError() {
	// Arrange.
	c := messagesrepo.Cursor{PageSize: 10, LastCreatedAt: time.Now()}
	cur, err := cursor.Encode(c)
	s.Require().NoError(err)
	req := getchathistory.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID(), ChatID: types.NewChatID(), Cursor: cur}
	errProblemsRepo := errors.New("problems repo error")
	s.problemsRepo.EXPECT().GetUnresolvedProblem(s.Ctx, req.ChatID, req.ManagerID).Return(nil, errProblemsRepo)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Contains(err.Error(), errProblemsRepo.Error())
	s.Empty(resp.Messages)
}

func (s *UseCaseSuite) TestHandle_MessagesRepoError() {
	// Arrange.
	req := getchathistory.Request{ID: types.NewRequestID(), ManagerID: types.NewUserID(), ChatID: types.NewChatID(), PageSize: 10}
	problem := &problemsrepo.Problem{
		ID:        types.NewProblemID(),
		ChatID:    req.ChatID,
		ManagerID: req.ManagerID,
		CreatedAt: time.Now(),
	}
	s.problemsRepo.EXPECT().GetUnresolvedProblem(s.Ctx, req.ChatID, req.ManagerID).Return(problem, nil)

	errMessagesRepo := errors.New("messages repo error")
	s.messagesRepo.EXPECT().
		GetProblemMessages(s.Ctx, problem.ID, req.PageSize, (*messagesrepo.Cursor)(nil)).
		Return(nil, nil, errMessagesRepo)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Contains(err.Error(), errMessagesRepo.Error())
	s.Empty(resp.Messages)
}
