package clientevents

import (
	"fmt"

	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	websocketstream "github.com/ekhvalov/bank-chat-service/internal/websocket-stream"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	if err := ev.Validate(); err != nil {
		return nil, fmt.Errorf("event validate: %v", err)
	}
	switch t := ev.(type) {
	case *eventstream.NewMessageEvent:
		return &NewMessageEvent{
			AuthorID:  pointer.PtrWithZeroAsNil(t.UserID),
			Body:      t.MessageBody,
			CreatedAt: &t.Time,
			ID:        t.ID,
			EventType: EventTypeNewMessageEvent,
			IsService: &t.IsService,
			MessageID: t.MessageID,
			RequestID: t.RequestID,
		}, nil
	case *eventstream.MessageSentEvent:
		return &MessageSentEvent{
			ID:        t.ID,
			EventType: EventTypeMessageSentEvent,
			MessageID: t.MessageID,
			RequestID: t.RequestID,
		}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", t)
	}
}
