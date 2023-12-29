//go:build integration

package problemsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	storeproblem "github.com/ekhvalov/bank-chat-service/internal/store/gen/problem"
	"github.com/ekhvalov/bank-chat-service/internal/testingh"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

type ProblemsRepoSuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsRepoSuite{DBSuite: testingh.NewDBSuite("TestProblemsRepoSuite")})
}

func (s *ProblemsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsRepoSuite) Test_CreateIfNotExists() {
	s.Run("problem does not exist, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)

		problem, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problemID)
		s.Require().NoError(err)
		s.Equal(problemID, problem.ID)
		s.Equal(chat.ID, problem.ChatID)
	})

	s.Run("resolved problem already exists, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().
			SetChatID(chat.ID).
			SetManagerID(types.NewUserID()).
			SetResolvedAt(time.Now()).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.NotEqual(problem.ID, problemID)
	})

	s.Run("problem already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.Equal(problem.ID, problemID)
	})
}

func (s *ProblemsRepoSuite) Test_CountManagerOpenProblems() {
	s.Run("manager has no open problems", func() {
		managerID := types.NewUserID()

		count, err := s.repo.CountManagerOpenProblems(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Empty(count)
	})

	s.Run("manager has open problems", func() {
		const (
			problemsCount         = 20
			resolvedProblemsCount = 3
		)

		managerID := types.NewUserID()
		problems := make([]types.ProblemID, 0, problemsCount)

		for i := 0; i < problemsCount; i++ {
			_, pID := s.createChatWithProblemAssignedTo(managerID)
			problems = append(problems, pID)
		}

		// Create problems for other managers.
		for i := 0; i < problemsCount; i++ {
			s.createChatWithProblemAssignedTo(types.NewUserID())
		}

		count, err := s.repo.CountManagerOpenProblems(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount, count)

		// Resolve some problems.
		for i := 0; i < resolvedProblemsCount; i++ {
			pID := problems[i*resolvedProblemsCount]
			_, err := s.Database.Problem(s.Ctx).
				Update().
				Where(storeproblem.ID(pID)).
				SetResolvedAt(time.Now()).
				Save(s.Ctx)
			s.Require().NoError(err)
		}

		count, err = s.repo.CountManagerOpenProblems(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount-resolvedProblemsCount, count)
	})
}

func (s *ProblemsRepoSuite) TestResolveAssignedProblem() {
	s.Run("no problem at all", func() {
		managerID := types.NewUserID()
		chatID := s.createChat()

		p, err := s.repo.ResolveAssignedProblem(s.Ctx, chatID, managerID)

		s.Require().ErrorIs(err, problemsrepo.ErrProblemNotFound)
		s.Nil(p)
	})

	s.Run("problem already resolved", func() {
		managerID := types.NewUserID()
		chatID := s.createChat()
		_, err := s.Database.Problem(s.Ctx).
			Create().
			SetChatID(chatID).
			SetManagerID(managerID).
			SetResolvedAt(time.Now()).
			Save(s.Ctx)
		s.Require().NoError(err)

		p, err := s.repo.ResolveAssignedProblem(s.Ctx, chatID, managerID)

		s.Require().ErrorIs(err, problemsrepo.ErrProblemNotFound)
		s.Nil(p)
	})

	s.Run("problem assigned to other manager", func() {
		managerID := types.NewUserID()
		chatID, _ := s.createChatWithProblemAssignedTo(types.NewUserID())

		p, err := s.repo.ResolveAssignedProblem(s.Ctx, chatID, managerID)

		s.Require().ErrorIs(err, problemsrepo.ErrProblemNotFound)
		s.Nil(p)
	})

	s.Run("problem successfully resolved", func() {
		managerID := types.NewUserID()
		chatID, _ := s.createChatWithProblemAssignedTo(managerID)

		p, err := s.repo.ResolveAssignedProblem(s.Ctx, chatID, managerID)

		s.Require().NoError(err)
		s.NotNil(p.ResolvedAt)
	})
}

func (s *ProblemsRepoSuite) createChatWithProblemAssignedTo(managerID types.UserID) (types.ChatID, types.ProblemID) {
	s.T().Helper()

	// 1 chat can have only 1 open problem.

	chatID := s.createChat()

	p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chatID).SetManagerID(managerID).Save(s.Ctx)
	s.Require().NoError(err)

	return chatID, p.ID
}

func (s *ProblemsRepoSuite) createChat() types.ChatID {
	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	return chat.ID
}
