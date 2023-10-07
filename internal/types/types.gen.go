package types

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrParse     = errors.New("parse error")
	ErrMarshal   = errors.New("marshal error")
	ErrUnmarshal = errors.New("unmarshal error")
	ErrInvalid   = errors.New("invalid value")
)

func Parse[T ChatID | MessageID](id string) (T, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return T(uuid.Nil), fmt.Errorf("%w: %v", ErrParse, err)
	}
	return T(uid), nil
}

func MustParse[T ChatID | MessageID](id string) T {
	uid, err := uuid.Parse(id)
	if err != nil {
		panic(fmt.Errorf("%w: %v", ErrParse, err))
	}
	return T(uid)
}

type ChatID uuid.UUID

var ChatIDNil = ChatID(uuid.Nil)

func NewChatID() ChatID {
	return ChatID(uuid.New())
}

func (c ChatID) MarshalText() (text []byte, err error) {
	text, err = uuid.UUID(c).MarshalText()
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrMarshal, err)
		return nil, err
	}
	return text, nil
}

func (c *ChatID) UnmarshalText(text []byte) error {
	uid, err := uuid.ParseBytes(text)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshal, err)
	}
	*c = ChatID(uid)
	return nil
}

func (c ChatID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c ChatID) Scan(src any) error {
	u := uuid.UUID(c)
	return u.Scan(src)
}

func (c ChatID) Validate() error {
	if c.IsZero() {
		return fmt.Errorf("%w: %v", ErrInvalid, c)
	}
	return nil
}

func (c ChatID) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c ChatID) String() string {
	return uuid.UUID(c).String()
}

func (c ChatID) IsZero() bool {
	return c == ChatIDNil
}

type MessageID uuid.UUID

var MessageIDNil = MessageID(uuid.Nil)

func NewMessageID() MessageID {
	return MessageID(uuid.New())
}

func (c MessageID) MarshalText() (text []byte, err error) {
	text, err = uuid.UUID(c).MarshalText()
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrMarshal, err)
		return nil, err
	}
	return text, nil
}

func (c *MessageID) UnmarshalText(text []byte) error {
	uid, err := uuid.ParseBytes(text)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshal, err)
	}
	*c = MessageID(uid)
	return nil
}

func (c MessageID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c MessageID) Scan(src any) error {
	u := uuid.UUID(c)
	return u.Scan(src)
}

func (c MessageID) Validate() error {
	if c.IsZero() {
		return fmt.Errorf("%w: %v", ErrInvalid, c)
	}
	return nil
}

func (c MessageID) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c MessageID) String() string {
	return uuid.UUID(c).String()
}

func (c MessageID) IsZero() bool {
	return c == MessageIDNil
}
