package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/ekhvalov/bank-chat-service/internal/store/gen/chat"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/message"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/predicate"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
)

const (
	pageSizeLimitMin = 10
	pageSizeLimitMax = 100
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	var limit int
	predicates := []predicate.Message{message.IsVisibleForClient(true)}
	if cursor != nil {
		if err := validateCursor(*cursor); err != nil {
			return nil, nil, err
		}
		limit = cursor.PageSize
		predicates = append(predicates, message.CreatedAtLT(cursor.LastCreatedAt))
	} else {
		if err := validatePageSize(pageSize); err != nil {
			return nil, nil, err
		}
		limit = pageSize
	}
	result, err := r.db.Chat(ctx).
		Query().
		Where(chat.ClientID(clientID)).
		QueryMessages().
		Order(message.ByCreatedAt(sql.OrderDesc())).
		Where(message.And(predicates...)).
		Limit(limit + 1).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	var cur *Cursor
	if len(result) > limit && result[limit-1] != nil {
		cur = &Cursor{PageSize: limit, LastCreatedAt: result[limit-1].CreatedAt}
	}
	messages := make([]Message, 0, limit)
	for i := 0; i < limit && i < len(result); i++ {
		if result[i] != nil {
			messages = append(messages, adaptStoreMessage(result[i]))
		}
	}
	return messages, cur, nil
}

func validatePageSize(pageSize int) error {
	if pageSize < pageSizeLimitMin || pageSize > pageSizeLimitMax {
		return ErrInvalidPageSize
	}
	return nil
}

func validateCursor(cursor Cursor) error {
	if err := validatePageSize(cursor.PageSize); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidCursor, ErrInvalidPageSize)
	}
	zero := time.Time{}
	if cursor.LastCreatedAt.Sub(zero).Nanoseconds() == 0 {
		return fmt.Errorf("%w: empty 'last created at'", ErrInvalidCursor)
	}
	return nil
}
