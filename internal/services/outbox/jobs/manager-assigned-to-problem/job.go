package managerassignedtoproblemjob

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -typed -package=managerassignedtoproblemjobmocks

const Name = "manager-assigned-to-problem"

type problemRepository interface {
	GetProblemByID(ctx context.Context, problemID types.ProblemID) (*problemsrepo.Problem, error)
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	messageRepo messageRepository  `option:"mandatory" validate:"required"`
	problemRepo problemRepository  `option:"mandatory" validate:"required"`
	managerLoad managerLoadService `option:"mandatory" validate:"required"`
	eventStream eventStream        `option:"mandatory" validate:"required"`
	log         *zap.Logger        `option:"mandatory" validate:"required"`
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Job{Options: opts}, nil
}

type Job struct {
	Options
	outbox.DefaultJob
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) error {
	j.log.Debug("got payload", zap.String("payload", payload))
	pl, err := Unmarshal(payload)
	if err != nil {
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	message, err := j.messageRepo.GetMessageByID(ctx, pl.MessageID)
	if err != nil {
		return fmt.Errorf("get message: %v", err)
	}

	canManagerTakeProblem, err := j.managerLoad.CanManagerTakeProblem(ctx, pl.ManagerID)
	if err != nil {
		return fmt.Errorf("can manager take problem: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		newChatEvent := eventstream.NewNewChatEvent(
			types.NewEventID(),
			message.ChatID,
			pl.ClientID,
			message.RequestID,
			canManagerTakeProblem,
		)

		if err := j.eventStream.Publish(ctx, pl.ManagerID, newChatEvent); err != nil {
			return fmt.Errorf("publish NewChatEvent: %v", err)
		}

		return nil
	})

	eg.Go(func() error {
		newMessageEvent := eventstream.NewNewMessageEvent(
			types.NewEventID(),
			message.RequestID,
			message.ChatID,
			message.ID,
			pl.ClientID,
			message.CreatedAt,
			message.Body,
			message.IsService,
		)

		if err := j.eventStream.Publish(ctx, pl.ClientID, newMessageEvent); err != nil {
			return fmt.Errorf("publish NewMessageEvent: %v", err)
		}

		return nil
	})

	return eg.Wait()
}
