package websocketstream

import (
	"encoding/json"
	"fmt"
	"io"

	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
)

// EventAdapter converts the event from the stream to the appropriate object.
type EventAdapter interface {
	Adapt(event eventstream.Event) (any, error)
}

// EventWriter write adapted event it to the socket.
type EventWriter interface {
	Write(event any, out io.Writer) error
}

type JSONEventWriter struct{}

func (JSONEventWriter) Write(event any, out io.Writer) error {
	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshall: %v", err)
	}
	if _, err := out.Write(b); err != nil {
		return fmt.Errorf("write: %v", err)
	}
	return nil
}
