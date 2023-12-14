// Package managerevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package managerevents

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
)

// Defines values for EventType.
const (
	EventTypeNewChatEvent EventType = "NewChatEvent"
)

// ChatID defines model for ChatID.
type ChatID = types.ChatID

// ClientID defines model for ClientID.
type ClientID = types.UserID

// Event defines model for Event.
type Event struct {
	union json.RawMessage
}

// ID defines model for EventID.
type ID = types.EventID

// EventType defines model for EventType.
type EventType string

// NewChatEvent defines model for NewChatEvent.
type NewChatEvent struct {
	CanTakeMoreProblems bool      `json:"canTakeMoreProblems"`
	ChatID              ChatID    `json:"chatId"`
	ClientID            ClientID  `json:"clientId"`
	ID                  ID        `json:"eventId"`
	EventType           EventType `json:"eventType"`
	RequestID           RequestID `json:"requestId"`
}

// RequestID defines model for RequestID.
type RequestID = types.RequestID

// AsNewChatEvent returns the union data inside the Event as a NewChatEvent
func (t Event) AsNewChatEvent() (NewChatEvent, error) {
	var body NewChatEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromNewChatEvent overwrites any union data inside the Event as the provided NewChatEvent
func (t *Event) FromNewChatEvent(v NewChatEvent) error {
	v.EventType = "NewChatEvent"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeNewChatEvent performs a merge with any union data inside the Event, using the provided NewChatEvent
func (t *Event) MergeNewChatEvent(v NewChatEvent) error {
	v.EventType = "NewChatEvent"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t Event) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"eventType"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t Event) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "NewChatEvent":
		return t.AsNewChatEvent()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t Event) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Event) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RVwW7bMAz9FYMbsIsTp92l0HFtMRRD22HtTkUPisPYWmxKk2hnReB/HyjbS9NlSzus",
	"pxgSH9/T45OygdzWzhISB1AbCHmJtY6fp6XmizP5WlpfawYFTWMWkAI/OAQFgb2hAlL4MSnshHQtiwNq",
	"WBwq5SdMf9+amNpZz8LhNJegoDBcNvNpbusMV2WrK9tmc02rSV5qngT0rckxM8ToSVdZbAxdl8JpZZD+",
	"Re+I26f4a0D/aorPW6TYaGFC7k1tSLP1sbO3Dj0/XPUSUQpvRVeXgiW8XoK628Bbj0tQ8CbbTjAbxpdd",
	"4Vrc7im6+5Ht5fb8wZix2ys6Ew+sNoDU1KDuYOdI909Vd+luwdZFgzHNuaZbvcJL6/Gzt/MK67g8dJlb",
	"W6EmaSPiLhay9zeDhyxLfR+gw4gxaF3aT/QwZHR5RIyWHMSMYfH4vcHwDKYvQ+FZ9F9gxuNCbB+lPpbw",
	"y6RHp39Mlu51W1K45XlhDrfAfXHcu/vfAimWGFramBfDlfB+0LRKbhonDImEIbnUpAv0SRxAgBRa9MFY",
	"AgXtUby5Dkk7AwreT4+mM0ijrJjCLHAzl48C+xcB5UVw3MMvOGkChmRpfVIgoddsqEjiPMI0ueYS/doE",
	"TAwnC4uB3vEUIp9UWpLxw0fkGyGROQVnKfTX4ng2i7fDEg/XRjtXmTwCs29BBIz/Cs+KHkS7dg9w/Sm6",
	"2KUgHqMP8f3arTnDFivraiRO+ipIofEVKFgHlWWVzXVV2sDqZHYyy9YBuvRpj/Pj86fYWHjf/QwAAP//",
	"/E2xtucGAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
