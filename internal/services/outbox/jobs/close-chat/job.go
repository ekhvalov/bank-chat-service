package closechatjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -typed -package=managerassignedtoproblemjobmocks

const Name = "close-chat"

type problemsRepository interface {
	ResolveAssignedProblem(ctx context.Context, chatID types.ChatID, managerID types.UserID) (*problemsrepo.Problem, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

type messagesRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	messageRepo messagesRepository `option:"mandatory" validate:"required"`
	eventStream eventStream        `option:"mandatory" validate:"required"`
	log         *zap.Logger        `option:"mandatory" validate:"required"`
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	opts.log = opts.log.Named(Name)
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
	pl, err := Unmarshal(payload)
	if err != nil {
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	msg, err := j.messageRepo.GetMessageByID(ctx, pl.MessageID)
	if err != nil {
		j.log.Error("get message by id", zap.Error(err))
		return fmt.Errorf("get message by id: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return j.publishChatClosedEvent(ctx, msg.ChatID, pl.ManagerID, pl.RequestID)
	})
	eg.Go(func() error {
		return j.publishNewMessageEvent(ctx, msg, pl.ClientID)
	})

	return eg.Wait()
}

func (j *Job) publishChatClosedEvent(
	ctx context.Context,
	chatID types.ChatID,
	managerID types.UserID,
	requestID types.RequestID,
) error {
	eventID := types.NewEventID()
	event := eventstream.NewChatClosedEvent(eventID, chatID, requestID, true)

	err := j.eventStream.Publish(ctx, managerID, event)
	if err != nil {
		j.log.Error("publish ChatClosedEvent", zap.Error(err))
		return fmt.Errorf("publish ChatClosedEvent: %v", err)
	}

	j.log.Debug(
		"published ChatClosedEvent",
		zap.Stringer("event_id", eventID),
		zap.Stringer("manager_id", managerID),
		zap.Stringer("chat_id", chatID),
	)
	return nil
}

func (j *Job) publishNewMessageEvent(
	ctx context.Context,
	msg *messagesrepo.Message,
	clientID types.UserID,
) error {
	eventID := types.NewEventID()
	event := eventstream.NewNewMessageEvent(
		eventID,
		msg.RequestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		msg.IsService,
	)

	err := j.eventStream.Publish(ctx, clientID, event)
	if err != nil {
		j.log.Error("publish NewMessageEvent", zap.Error(err))
		return fmt.Errorf("publish NewMessageEvent: %v", err)
	}

	j.log.Debug(
		"published NewMessageEvent",
		zap.Stringer("event_id", eventID),
		zap.Stringer("client_id", clientID),
		zap.Stringer("msg_id", msg.ID),
	)
	return nil
}
