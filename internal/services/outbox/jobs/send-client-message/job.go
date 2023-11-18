package sendclientmessagejob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	msgproducer "github.com/ekhvalov/bank-chat-service/internal/services/msg-producer"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

const Name = "send-client-message"

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	producer messageProducer   `option:"mandatory" validate:"required"`
	repo     messageRepository `option:"mandatory" validate:"required"`
	evStream eventStream       `option:"mandatory" validate:"required"`
	log      *zap.Logger       `option:"mandatory" validate:"required"`
}

type Job struct {
	Options
	outbox.DefaultJob
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Job{Options: opts}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) error {
	msgID, err := UnmarshalPayload(payload)
	if err != nil {
		return fmt.Errorf("unvarshal payload: %v", err)
	}

	msg, err := j.repo.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("get message: %v", err)
	}

	err = j.producer.ProduceMessage(ctx, msgproducer.Message{
		ID:         msg.ID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: true,
	})
	if err != nil {
		j.log.Error("produce message", zap.Error(err), zap.String("id", msgID.String()))
		return fmt.Errorf("produce message: %v", err)
	}
	requestID := types.NewRequestID()
	eventID := types.NewEventID()
	event := eventstream.NewNewMessageEvent(
		eventID,
		requestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		msg.IsService,
	)
	err = j.evStream.Publish(ctx, msg.AuthorID, event)
	if err != nil {
		j.log.Error("publish new message event", zap.Error(err), zap.String("id", msgID.String()))
		return fmt.Errorf("publish new message event: %v", err)
	}
	j.log.Debug(
		"NewMessageEvent published",
		zap.String("event_id", eventID.String()),
		zap.String("author_id", msg.AuthorID.String()),
	)

	j.log.Debug("message produced", zap.String("id", msgID.String()))
	return nil
}
