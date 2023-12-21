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
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/problem"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
	errEmptyCreatedAt  = errors.New("empty 'last created at'")
	zeroTime           = time.Time{}
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
	chatPredicates := []predicate.Chat{chat.ClientID(clientID)}
	messagePredicates := []predicate.Message{message.IsVisibleForClient(true)}

	messages, cursor, err := r.getChatMessages(ctx, chatPredicates, messagePredicates, pageSize, cursor)
	if err != nil {
		return nil, nil, fmt.Errorf("get client chat messages: %w", err)
	}

	return messages, cursor, nil
}

// GetProblemMessages returns Nth page of messages in the chat for manager side (specific problem).
func (r *Repo) GetProblemMessages(
	ctx context.Context,
	problemID types.ProblemID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	chatPredicates := []predicate.Chat{chat.HasProblemsWith(problem.ID(problemID))}
	messagePredicates := []predicate.Message{
		message.IsVisibleForManager(true),
		message.ProblemID(problemID),
	}

	messages, cursor, err := r.getChatMessages(ctx, chatPredicates, messagePredicates, pageSize, cursor)
	if err != nil {
		return nil, nil, fmt.Errorf("get problem messages: %w", err)
	}

	return messages, cursor, nil
}

func (r *Repo) getChatMessages(
	ctx context.Context,
	chatPredicates []predicate.Chat,
	messagePredicates []predicate.Message,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	var limit int
	if cursor != nil {
		if err := validateCursor(*cursor); err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
		limit = cursor.PageSize
		messagePredicates = append(messagePredicates, message.CreatedAtLT(cursor.LastCreatedAt))
	} else {
		if err := validatePageSize(pageSize); err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrInvalidPageSize, err)
		}
		limit = pageSize
	}

	result, err := r.db.Chat(ctx).
		Query().
		Where(chatPredicates...).
		QueryMessages().
		Order(message.ByCreatedAt(sql.OrderDesc())).
		Where(message.And(messagePredicates...)).
		Limit(limit + 1).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("query chat messages: %v", err)
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
		return err
	}
	if cursor.LastCreatedAt.Sub(zeroTime).Nanoseconds() == 0 {
		return errEmptyCreatedAt
	}
	return nil
}
