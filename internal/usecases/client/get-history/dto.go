package gethistory

import (
	"errors"
	"time"

	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/internal/validator"
)

type Request struct {
	ID       types.RequestID `validate:"required"`
	ClientID types.UserID    `validate:"required"`
	PageSize int             `validate:"omitempty,gte=10,lte=100"`
	Cursor   string          `validate:"omitempty,base64url"`
}

func (r Request) Validate() error {
	if r.PageSize <= 0 && r.Cursor == "" {
		return errors.New("neither cursor nor pageSize specified")
	} else if r.PageSize != 0 && r.Cursor != "" {
		return errors.New("cursor and pageSize specified")
	}
	return validator.Validator.Struct(r)
}

type Response struct {
	Messages   []Message
	NextCursor string
}

type Message struct {
	ID         types.MessageID
	AuthorID   types.UserID
	Body       string
	CreatedAt  time.Time
	IsReceived bool
	IsBlocked  bool
	IsService  bool
}
