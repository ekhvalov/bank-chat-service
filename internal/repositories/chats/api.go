package chatsrepo

import (
	"context"
	"fmt"
	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"

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
		Order(chat.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query chats: %v", err)
	}

	return adaptStoreChats(chats), nil
}

func (r *Repo) GetClientIDByChatID(ctx context.Context, chatID types.ChatID) (types.UserID, error) {
	ch, err := r.db.Chat(ctx).Get(ctx, chatID)
	if err != nil {
		if store.IsNotFound(err) {
			return types.UserIDNil, ErrChatsNotFound
		}
		return types.UserIDNil, fmt.Errorf("get chat: %v", err)
	}

	return ch.ClientID, nil
}
