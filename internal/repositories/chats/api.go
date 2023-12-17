package chatsrepo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/ekhvalov/bank-chat-service/internal/store/gen/chat"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/problem"
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

func (r *Repo) GetOpenProblemChatsForManager(ctx context.Context, managerID types.UserID) ([]Chat, error) {
	chats, err := r.db.Chat(ctx).Query().
		Where(
			chat.HasProblemsWith(
				problem.ManagerID(managerID),
				problem.ResolvedAtIsNil(),
			),
		).
		Order(chat.ByCreatedAt(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query chats: %v", err)
	}
	// if len(chats) == 0 {
	// 	return nil, ErrChatsNotFound
	// }

	return adaptStoreChats(chats), nil
}
