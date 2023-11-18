package sendclientmessagejob_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	msgproducer "github.com/ekhvalov/bank-chat-service/internal/services/msg-producer"
	sendclientmessagejob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-client-message"
	sendclientmessagejobmocks "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-client-message/mocks"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendclientmessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendclientmessagejobmocks.NewMockmessageRepository(ctrl)
	evStream := sendclientmessagejobmocks.NewMockeventStream(ctrl)
	job, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(msgProducer, msgRepo, evStream, zap.NewNop()))
	require.NoError(t, err)

	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            clientID,
		Body:                body,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           false,
	}
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	evStream.EXPECT().Publish(ctx, clientID, &newMessageEventMatcher{&eventstream.NewMessageEvent{
		ID:          types.EventID{},
		RequestID:   types.RequestID{},
		ChatID:      msg.ChatID,
		MessageID:   msg.ID,
		UserID:      clientID,
		Time:        msg.CreatedAt,
		MessageBody: msg.Body,
		IsService:   msg.IsService,
	}})

	msgProducer.EXPECT().ProduceMessage(gomock.Any(), msgproducer.Message{
		ID:         msgID,
		ChatID:     chatID,
		Body:       body,
		FromClient: true,
	}).Return(nil)

	// Action & assert.
	payload, err := sendclientmessagejob.MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}

type newMessageEventMatcher struct {
	*eventstream.NewMessageEvent
}

func (e *newMessageEventMatcher) Matches(x any) bool {
	ev, ok := x.(*eventstream.NewMessageEvent)
	if !ok {
		return false
	}
	e.ID = ev.ID
	e.RequestID = ev.RequestID
	return reflect.DeepEqual(e.NewMessageEvent, ev)
}

func (e *newMessageEventMatcher) String() string {
	return fmt.Sprintf("matches event: %v", e.NewMessageEvent)
}
