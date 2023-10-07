package types

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrParse     = errors.New("parse error")
	ErrScan      = errors.New("scan error")
	ErrMarshal   = errors.New("marshal error")
	ErrUnmarshal = errors.New("unmarshal error")
	ErrInvalid   = errors.New("invalid value")
)

func Parse[T ChatID | MessageID | ProblemID | UserID](id string) (T, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return T(uuid.Nil), fmt.Errorf("%w: %v", ErrParse, err)
	}
	return T(uid), nil
}

func MustParse[T ChatID | MessageID | ProblemID | UserID](id string) T {
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

func (c *ChatID) Scan(src any) error {
	u := uuid.UUID(*c)
	err := u.Scan(src)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrScan, err)
	}
	*c = ChatID(u)
	return nil
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

func (c *MessageID) Scan(src any) error {
	u := uuid.UUID(*c)
	err := u.Scan(src)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrScan, err)
	}
	*c = MessageID(u)
	return nil
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

type ProblemID uuid.UUID

var ProblemIDNil = ProblemID(uuid.Nil)

func NewProblemID() ProblemID {
	return ProblemID(uuid.New())
}

func (c ProblemID) MarshalText() (text []byte, err error) {
	text, err = uuid.UUID(c).MarshalText()
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrMarshal, err)
		return nil, err
	}
	return text, nil
}

func (c *ProblemID) UnmarshalText(text []byte) error {
	uid, err := uuid.ParseBytes(text)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshal, err)
	}
	*c = ProblemID(uid)
	return nil
}

func (c ProblemID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *ProblemID) Scan(src any) error {
	u := uuid.UUID(*c)
	err := u.Scan(src)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrScan, err)
	}
	*c = ProblemID(u)
	return nil
}

func (c ProblemID) Validate() error {
	if c.IsZero() {
		return fmt.Errorf("%w: %v", ErrInvalid, c)
	}
	return nil
}

func (c ProblemID) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c ProblemID) String() string {
	return uuid.UUID(c).String()
}

func (c ProblemID) IsZero() bool {
	return c == ProblemIDNil
}

type UserID uuid.UUID

var UserIDNil = UserID(uuid.Nil)

func NewUserID() UserID {
	return UserID(uuid.New())
}

func (c UserID) MarshalText() (text []byte, err error) {
	text, err = uuid.UUID(c).MarshalText()
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrMarshal, err)
		return nil, err
	}
	return text, nil
}

func (c *UserID) UnmarshalText(text []byte) error {
	uid, err := uuid.ParseBytes(text)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshal, err)
	}
	*c = UserID(uid)
	return nil
}

func (c UserID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *UserID) Scan(src any) error {
	u := uuid.UUID(*c)
	err := u.Scan(src)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrScan, err)
	}
	*c = UserID(u)
	return nil
}

func (c UserID) Validate() error {
	if c.IsZero() {
		return fmt.Errorf("%w: %v", ErrInvalid, c)
	}
	return nil
}

func (c UserID) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c UserID) String() string {
	return uuid.UUID(c).String()
}

func (c UserID) IsZero() bool {
	return c == UserIDNil
}
