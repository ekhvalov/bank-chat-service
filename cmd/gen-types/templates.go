package main

import "html/template"

func headerTemplate() *template.Template {
	text := `package {{.package}}

import (
	"errors"
	"fmt"

	"database/sql/driver"
	"github.com/google/uuid"
)

var (
	ErrParse     = errors.New("parse error")
	ErrMarshal   = errors.New("marshal error")
	ErrUnmarshal = errors.New("unmarshal error")
	ErrInvalid   = errors.New("invalid value")
)

func Parse[T {{.types}}](id string) (T, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return T(uuid.Nil), fmt.Errorf("%w: %v", ErrParse, err)
	}
	return T(uid), nil
}

func MustParse[T {{.types}}](id string) T {
	uid, err := uuid.Parse(id)
	if err != nil {
		panic(fmt.Errorf("%w: %v", ErrParse, err))
	}
	return T(uid)
}
`
	return template.Must(template.New("header").Parse(text))
}

func bodyTemplate() *template.Template {
	text := `
type {{.}} uuid.UUID

var (
	{{.}}Nil = {{.}}(uuid.Nil)
)

func New{{.}}() {{.}} {
	return {{.}}(uuid.New())
}

func (c {{.}}) MarshalText() (text []byte, err error) {
	text, err = uuid.UUID(c).MarshalText()
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrMarshal, err)
		return nil, err
	}
	return text, nil
}

func (c *{{.}}) UnmarshalText(text []byte) error {
	uid, err := uuid.ParseBytes(text)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshal, err)
	}
	*c = {{.}}(uid)
	return nil
}

func (c {{.}}) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c {{.}}) Scan(src any) error {
	u := uuid.UUID(c)
	return u.Scan(src)
}

func (c {{.}}) Validate() error {
	if c.IsZero() {
		return fmt.Errorf("%w: %v", ErrInvalid, c)
	}
	return nil
}

func (c {{.}}) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c {{.}}) String() string {
	return uuid.UUID(c).String()
}

func (c {{.}}) IsZero() bool {
	return c == {{.}}Nil
}
`

	return template.Must(template.New("body").Parse(text))
}
