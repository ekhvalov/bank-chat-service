package problemsrepo

import (
	"context"
	"fmt"

	"github.com/ekhvalov/bank-chat-service/internal/store"
	"github.com/ekhvalov/bank-chat-service/internal/store/problem"
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
