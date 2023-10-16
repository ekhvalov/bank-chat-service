package messagesrepo

import (
	"fmt"
	"time"

	"go.uber.org/mock/gomock"
)

var _ gomock.Matcher = CursorMatcher{}

// CursorMatcher is intended to be used only in tests.
type CursorMatcher struct {
	c Cursor
}

func NewCursorMatcher(c Cursor) CursorMatcher {
	return CursorMatcher{c: c}
}

func (cm CursorMatcher) Matches(x interface{}) bool {
	v, ok := x.(*Cursor)
	if !ok {
		return false
	}

	return cm.c.PageSize == v.PageSize &&
		truncateToMilliseconds(cm.c.LastCreatedAt).Equal(truncateToMilliseconds(v.LastCreatedAt))
}

func (cm CursorMatcher) String() string {
	return fmt.Sprintf("{ps=%d, last_created_at=%d}", cm.c.PageSize, cm.c.LastCreatedAt.UnixNano())
}

func truncateToMilliseconds(t time.Time) time.Time {
	return t.Truncate(time.Millisecond)
}
