package sendmanagermessagejob_test

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
	sendmanagermessagejob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-manager-message"
	sendmanagermessagejobmocks "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-manager-message/mocks"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/pkg/jobpayload"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendmanagermessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendmanagermessagejobmocks.NewMockmessagesRepository(ctrl)
	evStream := sendmanagermessagejobmocks.NewMockeventStream(ctrl)
	chatsRepo := sendmanagermessagejobmocks.NewMockchatsRepository(ctrl)
	job, err := sendmanagermessagejob.New(
		sendmanagermessagejob.NewOptions(msgProducer, msgRepo, chatsRepo, evStream, zap.NewNop()),
	)
	require.NoError(t, err)

	managerID := types.NewUserID()
	clientID := types.NewUserID()
	chatID := types.NewChatID()

	msg := createManagerMessage(chatID, managerID)
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msg.ID).Return(msg, nil)

	chatsRepo.EXPECT().GetClientIDByChatID(gomock.Any(), chatID).Return(clientID, nil)

	// Notify manager
	evStream.EXPECT().Publish(gomock.Any(), managerID, createNewMessageEventMatcher(msg))
	// Notify client
	evStream.EXPECT().Publish(gomock.Any(), clientID, createNewMessageEventMatcher(msg))

	msgProducer.EXPECT().
		ProduceMessage(gomock.Any(), msgproducer.Message{ID: msg.ID, ChatID: chatID, Body: msg.Body, FromClient: false}).
		Return(nil)

	payload, err := jobpayload.Marshal(msg.ID)
	require.NoError(t, err)

	// Action
	err = job.Handle(ctx, payload)

	// Assert
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

func createManagerMessage(chatID types.ChatID, managerID types.UserID) *messagesrepo.Message {
	return &messagesrepo.Message{
		ID:                  types.NewMessageID(),
		ChatID:              chatID,
		AuthorID:            managerID,
		Body:                "Hello!",
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: true,
		IsBlocked:           false,
		IsService:           false,
	}
}

func createNewMessageEventMatcher(msg *messagesrepo.Message) *newMessageEventMatcher {
	return &newMessageEventMatcher{&eventstream.NewMessageEvent{
		ID:          types.EventID{},
		RequestID:   types.RequestID{},
		ChatID:      msg.ChatID,
		MessageID:   msg.ID,
		UserID:      msg.AuthorID,
		Time:        msg.CreatedAt,
		MessageBody: msg.Body,
		IsService:   msg.IsService,
	}}
}
