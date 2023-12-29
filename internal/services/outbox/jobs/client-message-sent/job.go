package clientmessagesentjob

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/pkg/jobpayload"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=clientmessagesentjobmocks

const Name = "client-message-sent"

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type problemsRepo interface {
	GetProblemByMessageID(ctx context.Context, messageID types.MessageID) (*problemsrepo.Problem, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	messagesRepo messageRepository `option:"mandatory" validate:"required"`
	problemsRepo problemsRepo      `option:"mandatory" validate:"required"`
	evStream     eventStream       `option:"mandatory" validate:"required"`
	log          *zap.Logger       `option:"mandatory" validate:"required"`
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
	msgID, err := jobpayload.Unmarshal(payload)
	if err != nil {
		return fmt.Errorf("unvarshal payload: %v", err)
	}

	msg, err := j.messagesRepo.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("get message: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	// MessageSentEvent
	eg.Go(func() error {
		eventID := types.NewEventID()
		event := eventstream.NewMessageSentEvent(eventID, types.NewRequestID(), msg.ID)
		if err := j.evStream.Publish(ctx, msg.AuthorID, event); err != nil {
			j.log.Error("publish MessageSentEvent", zap.Error(err), zap.Stringer("message_id", msgID))
			return fmt.Errorf("publish message sent event: %v", err)
		}
		j.log.Debug(
			" published",
			zap.Stringer("event_id", eventID),
			zap.String("event_type", "MessageSentEvent"),
			zap.Stringer("author_id", msg.AuthorID),
			zap.Stringer("msg_id", msg.ID),
		)
		return nil
	})

	// NewMessageEvent
	eg.Go(func() error {
		problem, err := j.problemsRepo.GetProblemByMessageID(ctx, msgID)
		if err != nil {
			if errors.Is(err, problemsrepo.ErrProblemNotFound) {
				return nil // Manager is not assigned, skip.
			}
			return fmt.Errorf("get problem by message ID: %v", err)
		}

		eventID := types.NewEventID()
		event := eventstream.NewNewMessageEvent(
			eventID,
			msg.RequestID,
			msg.ChatID,
			msg.ID,
			msg.AuthorID,
			msg.CreatedAt,
			msg.Body,
			false,
		)
		if err := j.evStream.Publish(ctx, problem.ManagerID, event); err != nil {
			j.log.Error("publish NewMessageEvent", zap.Error(err), zap.Stringer("message_id", msgID))
		}
		j.log.Debug(
			" published",
			zap.Stringer("event_id", eventID),
			zap.String("event_type", "NewMessageEvent"),
			zap.Stringer("author_id", msg.AuthorID),
			zap.Stringer("msg_id", msg.ID),
		)

		return nil
	})

	return eg.Wait()
}
