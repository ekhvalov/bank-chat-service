package messagesrepo

import (
	"context"
	"time"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func (r *Repo) MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error {
	return r.db.Message(ctx).
		UpdateOneID(msgID).
		SetIsVisibleForManager(true).
		SetCheckedAt(time.Now()).
		Exec(ctx)
}

func (r *Repo) BlockMessage(ctx context.Context, msgID types.MessageID) error {
	return r.db.Message(ctx).
		UpdateOneID(msgID).
		SetIsBlocked(true).
		SetIsVisibleForManager(false).
		SetCheckedAt(time.Now()).
		Exec(ctx)
}
