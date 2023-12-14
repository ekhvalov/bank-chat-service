package eventstream

import (
	"fmt"
	"reflect"
)

type NewMessageEventMatcher struct {
	*NewMessageEvent
}

func (e *NewMessageEventMatcher) Matches(x any) bool {
	ev, ok := x.(*NewMessageEvent)
	if !ok {
		return false
	}
	e.ID = ev.ID
	e.RequestID = ev.RequestID
	return reflect.DeepEqual(e.NewMessageEvent, ev)
}

func (e *NewMessageEventMatcher) String() string {
	return fmt.Sprintf("matches event: %v", e.NewMessageEvent)
}

type NewChatEventMatcher struct {
	*NewChatEvent
}

func (e *NewChatEventMatcher) Matches(x any) bool {
	ev, ok := x.(*NewChatEvent)
	if !ok {
		return false
	}
	e.ID = ev.ID
	e.RequestID = ev.RequestID
	return reflect.DeepEqual(e.NewChatEvent, ev)
}

func (e *NewChatEventMatcher) String() string {
	return fmt.Sprintf("matches event: %v", e.NewChatEvent)
}
