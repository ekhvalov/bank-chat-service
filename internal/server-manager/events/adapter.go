package managerevents

import (
	"fmt"

	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	websocketstream "github.com/ekhvalov/bank-chat-service/internal/websocket-stream"
)

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	if err := ev.Validate(); err != nil {
		return nil, fmt.Errorf("event validate: %v", err)
	}
	switch t := ev.(type) {
	case *eventstream.NewChatEvent:
		return NewChatEvent{
			ID:                  t.ID,
			ChatID:              t.ChatID,
			ClientID:            t.ClientID,
			EventType:           EventTypeNewChatEvent,
			RequestID:           t.RequestID,
			CanTakeMoreProblems: t.CanTakeMoreProblems,
		}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", t)
	}
}
