package chatsrepo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	id, err := r.db.Chat(ctx).
		Create().
		SetClientID(userID).
		OnConflict(
			sql.ConflictColumns("client_id"),
		).
		Ignore().
		ID(ctx)
	if err != nil {
		return types.ChatIDNil, fmt.Errorf("create chat: %v", err)
	}
	return id, nil
}
