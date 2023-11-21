package problemsrepo

import (
	"context"
	"fmt"

	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/predicate"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/problem"
	"github.com/ekhvalov/bank-chat-service/internal/types"
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
