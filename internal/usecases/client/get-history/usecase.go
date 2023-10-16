package gethistory

import (
	"context"
	"errors"
	"fmt"

	"github.com/ekhvalov/bank-chat-service/internal/cursor"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=gethistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetClientChatMessages(
		ctx context.Context,
		clientID types.UserID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("options validate error: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}
	var cur *messagesrepo.Cursor
	if req.Cursor != "" {
		if err := cursor.Decode(req.Cursor, &cur); err != nil {
			return Response{}, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
	}
	messages, nextCur, err := u.msgRepo.GetClientChatMessages(ctx, req.ClientID, req.PageSize, cur)
	if err != nil {
		if errors.Is(err, messagesrepo.ErrInvalidCursor) {
			return Response{}, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
		return Response{}, fmt.Errorf("GetClientChatMessages: %v", err)
	}
	resp := Response{}
	if nextCur != nil {
		if resp.NextCursor, err = cursor.Encode(nextCur); err != nil {
			return Response{}, fmt.Errorf("cursor encode: %v", err)
		}
	}
	resp.Messages = make([]Message, len(messages))
	for i, m := range messages {
		resp.Messages[i] = Message{
			ID:         m.ID,
			AuthorID:   m.AuthorID,
			Body:       m.Body,
			CreatedAt:  m.CreatedAt,
			IsReceived: m.IsVisibleForManager && !m.IsBlocked,
			IsBlocked:  m.IsBlocked,
			IsService:  m.IsService,
		}
	}
	return resp, nil
}
