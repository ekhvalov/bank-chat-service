package problemsrepo

import (
	"context"
	"fmt"
	"time"

	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/chat"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/message"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/predicate"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/problem"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	p, err := r.db.Problem(ctx).Query().Where(problem.ChatID(chatID), problem.ResolvedAtIsNil()).First(ctx)
	if nil == err {
		return p.ID, nil
	}

	if !store.IsNotFound(err) {
		return types.ProblemIDNil, fmt.Errorf("query problem: %v", err)
	}

	p, err = r.db.Problem(ctx).Create().SetChatID(chatID).Save(ctx)
	if err != nil {
		return types.ProblemIDNil, fmt.Errorf("create problem: %v", err)
	}
	return p.ID, nil
}

func (r *Repo) CountManagerOpenProblems(ctx context.Context, managerID types.UserID) (int, error) {
	conditions := []predicate.Problem{problem.ManagerID(managerID), problem.ResolvedAtIsNil()}
	count, err := r.db.Problem(ctx).Query().Where(conditions...).Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count: %v", err)
	}
	return count, nil
}

func (r *Repo) GetProblemByMessageID(ctx context.Context, messageID types.MessageID) (*Problem, error) {
	p, err := r.db.Problem(ctx).Query().Where(problem.HasMessagesWith(message.ID(messageID))).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("get problem: %v", err)
	}
	if nil == p {
		return nil, ErrProblemNotFound
	}

	return pointer.Ptr(adaptStoreProblem(p)), nil
}

func (r *Repo) GetUnresolvedProblem(ctx context.Context, chatID types.ChatID, managerID types.UserID) (*Problem, error) {
	p, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.HasChatWith(chat.ID(chatID)),
			problem.ManagerID(managerID),
			problem.ResolvedAtIsNil(),
		).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("get problem: %v", err)
	}
	if nil == p {
		return nil, ErrProblemNotFound
	}

	return pointer.Ptr(adaptStoreProblem(p)), nil
}

func (r *Repo) GetAssignedProblemID(ctx context.Context, managerID types.UserID, chatID types.ChatID) (types.ProblemID, error) {
	problemID, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.ChatID(chatID),
			problem.ManagerID(managerID),
			problem.ResolvedAtIsNil(),
		).
		OnlyID(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return types.ProblemIDNil, ErrProblemNotFound
		}
		return types.ProblemIDNil, fmt.Errorf("find problem id: %v", err)
	}

	return problemID, nil
}

func (r *Repo) ResolveAssignedProblem(ctx context.Context, chatID types.ChatID, managerID types.UserID) (*Problem, error) {
	p, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.ChatID(chatID),
			problem.ManagerID(managerID),
			problem.ResolvedAtIsNil(),
		).
		First(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrProblemNotFound
		}
		return nil, fmt.Errorf("get assigned problem: %v", err)
	}

	p, err = r.db.Problem(ctx).UpdateOne(p).SetResolvedAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolve assigned problem: %v", err)
	}

	return pointer.Ptr(adaptStoreProblem(p)), nil
}
