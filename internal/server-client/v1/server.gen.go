// Package clientv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package clientv1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for ErrorCode.
const (
	ErrorCodeCreateChatError    ErrorCode = 1000
	ErrorCodeCreateProblemError ErrorCode = 1001
)

// Error defines model for Error.
type Error struct {
	// Code contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
	Code    ErrorCode `json:"code"`
	Details *string   `json:"details,omitempty"`
	Message string    `json:"message"`
}

// ErrorCode contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
type ErrorCode int

// GetHistoryRequest defines model for GetHistoryRequest.
type GetHistoryRequest struct {
	Cursor   *string `json:"cursor,omitempty"`
	PageSize *int    `json:"pageSize,omitempty"`
}

// GetHistoryResponse defines model for GetHistoryResponse.
type GetHistoryResponse struct {
	Data  *MessagesPage `json:"data,omitempty"`
	Error *Error        `json:"error,omitempty"`
}

// Message defines model for Message.
type Message struct {
	AuthorId   *types.UserID   `json:"authorId,omitempty"`
	Body       string          `json:"body"`
	CreatedAt  time.Time       `json:"createdAt"`
	Id         types.MessageID `json:"id"`
	IsBlocked  bool            `json:"isBlocked"`
	IsReceived bool            `json:"isReceived"`
	IsService  bool            `json:"isService"`
}

// MessageHeader defines model for MessageHeader.
type MessageHeader struct {
	AuthorID  *types.UserID   `json:"authorId,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
	ID        types.MessageID `json:"id"`
}

// MessagesPage defines model for MessagesPage.
type MessagesPage struct {
	Messages []Message `json:"messages"`
	Next     string    `json:"next"`
}

// SendMessageRequest defines model for SendMessageRequest.
type SendMessageRequest struct {
	MessageBody string `json:"messageBody"`
}

// SendMessageResponse defines model for SendMessageResponse.
type SendMessageResponse struct {
	Data  *MessageHeader `json:"data,omitempty"`
	Error *Error         `json:"error,omitempty"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostGetHistoryParams defines parameters for PostGetHistory.
type PostGetHistoryParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostSendMessageParams defines parameters for PostSendMessage.
type PostSendMessageParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetHistoryJSONRequestBody defines body for PostGetHistory for application/json ContentType.
type PostGetHistoryJSONRequestBody = GetHistoryRequest

// PostSendMessageJSONRequestBody defines body for PostSendMessage for application/json ContentType.
type PostSendMessageJSONRequestBody = SendMessageRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /getHistory)
	PostGetHistory(ctx echo.Context, params PostGetHistoryParams) error

	// (POST /sendMessage)
	PostSendMessage(ctx echo.Context, params PostSendMessageParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostGetHistory converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetHistory(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetHistoryParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostGetHistory(ctx, params)
	return err
}

// PostSendMessage converts echo context to params.
func (w *ServerInterfaceWrapper) PostSendMessage(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostSendMessageParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostSendMessage(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/getHistory", wrapper.PostGetHistory)
	router.POST(baseURL+"/sendMessage", wrapper.PostSendMessage)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xWT2/jthP9KsT8focWkCy56WEhoIf8aTcpukCwSdEFUh9oaWKxkUgtOXKTBvruxZC0",
	"JcdON902i55skRzOzHtvhvMIpWk7o1GTg+IROmlli4TWf314jx97dHRxdo6yQstrSkMBdfhMQMsWoYAP",
	"aTyZXpxBAhY/9spiBQXZHhNwZY2tZOtbY1tJUEDfqwoSoIeO7R1ZpVeQwH26Mmlc5B8324Yw3U1V2xlL",
	"IWKqoYCVorpfzkrTZnhXr2Vj1tlS6ru0rCWlDu1alZgpTWi1bDJ/NwzDMGyi8wl/b63xWXbWdGhJoV8u",
	"TYX8+3+Lt1DA/7IRtCxaZ970lA8OCVRIUjXedjfDIYEWnZMrPLA3TJG72R5Mgv/FkMDopHiECl1pVUfK",
	"MCWl0SSVduL8+vpSIB8UbOeE1JVwHZbqVpVi2Tul0TnRmJUqd859RTWKRjoSbe9ILFH82uf5EX4n5nme",
	"fz2DBFD3LRQ3/J3M83y+SKBVWrW8+m2eb/lknFdeIPcp26RraVkqjvPaJnFqURKe1pIC7snTrUtrlg22",
	"YZfzf4t0rhwZ+xBlcYCr3rrA4R7ynVzhlfrDg9fK+xD2nMPeJjHfz8FrZOrYdUY73PdcSZKfUsm7wKm7",
	"ZGKHBHAjuE9KK8TxbhTPrnPZU23sRfV5RfazQ/tKFZbA0lQPB/koPcnVMe0EXUnClFSLe5EPCajPTDDi",
	"9mo5KnfSmPIOq0miS2MalNqH7d5jiWr9/P5VuPvQ9pO24DP2mE4R3PExjWd6+WJU0NjQ/6GO4gtwHOzO",
	"vry8XldFMb1nEntdWR0ifsx2QmZoJ3tcxgfE/1eErXthd2KIYprSWvnA3xrv6cVPloNowDFeoa7ixc82",
	"7Wh3EjtFK+9/Qr1i7I7y2J83C/PkZTH4u/b8/wu9O1bO327ePGlg2VtFD1e8FxwvUVq0XD3j1w8bYf74",
	"yzXE+cR3BL87KrUm6oJMlL41nh1FDe+cSH0nrvqOdSj4gRWnjUJN4vjyAhJYo3VhaljPORHToZadggKO",
	"ZvnsCBIvXB9fttq+fB41E+jbnT3eIgmWsqjDSR4VGF3J+9xL4NI4Gt9Q72AcNG8OIzgeyfYG0WERSEdH",
	"G9HwAITaRye7rlGl95795jjEx8kM+lds7Q8YT8qQB1q/EJTkMfomz18lgChWH8Eu4Ju6F41yNIvqytyo",
	"9Oe54nIQGn8XsU4EGcGDH/N3mLdJAf13iTvQZb4wc4f6zPPUidjIA3mT3uBRnXaFmwVjxo/EBvPdC89w",
	"jY3pWi7vcAoS6G0TG0SRZY0pZVMbR8Wb/E2ecc0vhj8DAAD//9Kfh+X6DQAA",
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
