package clientmessageblockedjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/pkg/jobpayload"
)

const Name = "client-message-blocked"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=clientmessageblockedjobmocks

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	repo     messageRepository `option:"mandatory" validate:"required"`
	evStream eventStream       `option:"mandatory" validate:"required"`
	log      *zap.Logger       `option:"mandatory" validate:"required"`
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

	msg, err := j.repo.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("get message: %v", err)
	}

	eventID := types.NewEventID()
	event := eventstream.NewMessageBlockedEvent(eventID, types.NewRequestID(), msg.ID)
	err = j.evStream.Publish(ctx, msg.AuthorID, event)
	if err != nil {
		j.log.Error("publish message blocked event", zap.Error(err), zap.String("id", msgID.String()))
		return fmt.Errorf("publish message blocked event: %v", err)
	}
	j.log.Debug(
		"MessageBlockedEvent published",
		zap.Stringer("event_id", eventID),
		zap.Stringer("author_id", msg.AuthorID),
		zap.Stringer("msg_id", msg.ID),
	)

	return nil
}
