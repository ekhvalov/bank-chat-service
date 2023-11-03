package sendmessage

import (
	"context"
	"errors"
	"fmt"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrChatNotCreated    = errors.New("chat not created")
	ErrProblemNotCreated = errors.New("problem not created")
	ErrMessageNotCreated = errors.New("message not created")
)

type chatsRepository interface {
	CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error)
}

type messagesRepository interface {
	GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*messagesrepo.Message, error)
	CreateClientVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatsRepo    chatsRepository    `option:"mandatory" validate:"required"`
	messagesRepo messagesRepository `option:"mandatory" validate:"required"`
	problemsRepo problemsRepository `option:"mandatory" validate:"required"`
	tor          transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	return UseCase{Options: opts}, opts.Validate()
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	var message *messagesrepo.Message
	tx := func(ctx context.Context) error {
		var err error
		message, err = u.messagesRepo.GetMessageByRequestID(ctx, req.ID)
		if err != nil {
			if !errors.Is(err, messagesrepo.ErrMsgNotFound) {
				return fmt.Errorf("GetMessageByRequestID: %v", err)
			}
			var chatID types.ChatID
			chatID, err = u.chatsRepo.CreateIfNotExists(ctx, req.ClientID)
			if err != nil {
				return ErrChatNotCreated
			}
			var problemID types.ProblemID
			problemID, err = u.problemsRepo.CreateIfNotExists(ctx, chatID)
			if err != nil {
				return ErrProblemNotCreated
			}
			message, err = u.messagesRepo.CreateClientVisible(ctx, req.ID, problemID, chatID, req.ClientID, req.MessageBody)
			if err != nil {
				return ErrMessageNotCreated
			}
		}
		return nil
	}
	if err := u.tor.RunInTx(ctx, tx); err != nil {
		return Response{}, fmt.Errorf("tx commit: %w", err)
	}

	return Response{
		AuthorID:  req.ClientID,
		MessageID: message.ID,
		CreatedAt: message.CreatedAt,
	}, nil
}
