package messagesrepo

import (
	"context"
	"errors"
	"fmt"

	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/chat"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/message"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/problem"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

var ErrMsgNotFound = errors.New("message not found")

func (r *Repo) GetMessageByID(ctx context.Context, id types.MessageID) (*Message, error) {
	msg, err := r.db.Message(ctx).Get(ctx, id)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrMsgNotFound
		}
		return nil, fmt.Errorf("get message: %v", err)
	}
	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

func (r *Repo) GetInitialMessageByProblemID(ctx context.Context, problemID types.ProblemID) (*Message, error) {
	msg, err := r.db.Message(ctx).
		Query().
		Where(
			message.HasChatWith(chat.HasProblemsWith(problem.ID(problemID))),
			message.IsVisibleForManager(true),
		).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("find message: %v", err)
	}

	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

func (r *Repo) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*Message, error) {
	msg, err := r.db.Message(ctx).Query().Where(message.InitialRequestID(reqID)).First(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrMsgNotFound
		}
		return nil, fmt.Errorf("query by request id: %v", err)
	}
	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

// CreateClientVisible creates a message that is visible only to the client.
func (r *Repo) CreateClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).
		Create().
		SetInitialRequestID(reqID).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetBody(msgBody).
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(false).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create message: %v", err)
	}
	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

func (r *Repo) CreateServiceClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).
		Create().
		SetInitialRequestID(reqID).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetBody(msgBody).
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(false).
		SetIsService(true).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create message: %v", err)
	}
	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

func (r *Repo) CreateFullVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).
		Create().
		SetID(types.NewMessageID()).
		SetChatID(chatID).
		SetProblemID(problemID).
		SetInitialRequestID(reqID).
		SetAuthorID(authorID).
		SetBody(msgBody).
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(true).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create message: %v", err)
	}

	return pointer.Ptr(adaptStoreMessage(msg)), nil
}
