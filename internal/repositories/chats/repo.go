package chatsrepo

import (
	"context"
	"fmt"

	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate options-gen -out-filename=repo_options.gen.go -from-struct=Options
type Options struct {
	db *store.Database `option:"mandatory" validate:"required"`
}

type Repo struct {
	Options
}

func (r *Repo) GetClientIDByChatID(ctx context.Context, chatID types.ChatID) (types.UserID, error) {
	chat, err := r.db.Chat(ctx).Get(ctx, chatID)
	if err != nil {
		if store.IsNotFound(err) {
			return types.UserIDNil, ErrChatsNotFound
		}
		return types.UserIDNil, fmt.Errorf("get chat: %v", err)
	}

	return chat.ClientID, nil
}

func New(opts Options) (*Repo, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Repo{Options: opts}, nil
}
