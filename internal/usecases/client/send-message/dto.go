package sendmessage

import (
	"fmt"
	"time"

	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/internal/validator"
)

type Request struct {
	ID          types.RequestID `validate:"required"`
	ClientID    types.UserID    `validate:"required"`
	MessageBody string          `validate:"min=1,max=3000"`
}

func (r Request) Validate() error {
	if err := validator.Validator.Struct(r); err != nil {
		return fmt.Errorf("validate: %v", err)
	}
	return nil
}

type Response struct {
	AuthorID  types.UserID
	MessageID types.MessageID
	CreatedAt time.Time
}
