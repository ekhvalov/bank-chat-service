// Package apimanagerv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package apimanagerv1

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for ErrorCode.
const (
	ErrorCodeManagerOverloadedError ErrorCode = 5000
	ErrorCodeNoActiveProblemInChat  ErrorCode = 5001
)

// Chat defines model for Chat.
type Chat struct {
	ChatId   types.ChatID `json:"chatId"`
	ClientId types.UserID `json:"clientId"`
}

// ChatId defines model for ChatId.
type ChatId struct {
	ChatId types.ChatID `json:"chatId"`
}

// ChatList defines model for ChatList.
type ChatList struct {
	Chats []Chat `json:"chats"`
}

// CloseChatRequest defines model for CloseChatRequest.
type CloseChatRequest = ChatId

// CloseChatResponse defines model for CloseChatResponse.
type CloseChatResponse = NullDataResponse

// Error defines model for Error.
type Error struct {
	// Code contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
	Code    ErrorCode `json:"code"`
	Details *string   `json:"details,omitempty"`
	Message string    `json:"message"`
}

// ErrorCode contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
type ErrorCode int

// FreeHandsAvailability defines model for FreeHandsAvailability.
type FreeHandsAvailability struct {
	Available bool `json:"available"`
}

// FreeHandsResponse defines model for FreeHandsResponse.
type FreeHandsResponse = NullDataResponse

// GetChatHistoryRequest defines model for GetChatHistoryRequest.
type GetChatHistoryRequest struct {
	ChatId   types.ChatID `json:"chatId"`
	Cursor   *string      `json:"cursor,omitempty"`
	PageSize *int         `json:"pageSize,omitempty"`
}

// GetChatHistoryResponse defines model for GetChatHistoryResponse.
type GetChatHistoryResponse struct {
	Data  *MessagesPage `json:"data,omitempty"`
	Error *Error        `json:"error,omitempty"`
}

// GetChatsResponse defines model for GetChatsResponse.
type GetChatsResponse struct {
	Data  *ChatList `json:"data,omitempty"`
	Error *Error    `json:"error,omitempty"`
}

// GetFreeHandsBtnAvailabilityResponse defines model for GetFreeHandsBtnAvailabilityResponse.
type GetFreeHandsBtnAvailabilityResponse struct {
	Data  *FreeHandsAvailability `json:"data,omitempty"`
	Error *Error                 `json:"error,omitempty"`
}

// Message defines model for Message.
type Message struct {
	AuthorId   types.UserID    `json:"authorId"`
	Body       string          `json:"body"`
	CreatedAt  time.Time       `json:"createdAt"`
	Id         types.MessageID `json:"id"`
	IsReceived bool            `json:"isReceived"`
}

// MessageWithoutBody defines model for MessageWithoutBody.
type MessageWithoutBody struct {
	AuthorId  types.UserID    `json:"authorId"`
	CreatedAt time.Time       `json:"createdAt"`
	Id        types.MessageID `json:"id"`
}

// MessagesPage defines model for MessagesPage.
type MessagesPage struct {
	Messages []Message `json:"messages"`
	Next     string    `json:"next"`
}

// NullDataResponse defines model for NullDataResponse.
type NullDataResponse struct {
	Data  *interface{} `json:"data,omitempty"`
	Error *Error       `json:"error,omitempty"`
}

// SendMessageRequest defines model for SendMessageRequest.
type SendMessageRequest struct {
	ChatId      types.ChatID `json:"chatId"`
	MessageBody string       `json:"messageBody"`
}

// SendMessageResponse defines model for SendMessageResponse.
type SendMessageResponse struct {
	Data  *MessageWithoutBody `json:"data,omitempty"`
	Error *Error              `json:"error,omitempty"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostCloseChatParams defines parameters for PostCloseChat.
type PostCloseChatParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostFreeHandsParams defines parameters for PostFreeHands.
type PostFreeHandsParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetChatHistoryParams defines parameters for PostGetChatHistory.
type PostGetChatHistoryParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetChatsParams defines parameters for PostGetChats.
type PostGetChatsParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetFreeHandsBtnAvailabilityParams defines parameters for PostGetFreeHandsBtnAvailability.
type PostGetFreeHandsBtnAvailabilityParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostSendMessageParams defines parameters for PostSendMessage.
type PostSendMessageParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostCloseChatJSONRequestBody defines body for PostCloseChat for application/json ContentType.
type PostCloseChatJSONRequestBody = CloseChatRequest

// PostGetChatHistoryJSONRequestBody defines body for PostGetChatHistory for application/json ContentType.
type PostGetChatHistoryJSONRequestBody = GetChatHistoryRequest

// PostSendMessageJSONRequestBody defines body for PostSendMessage for application/json ContentType.
type PostSendMessageJSONRequestBody = SendMessageRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostCloseChatWithBody request with any body
	PostCloseChatWithBody(ctx context.Context, params *PostCloseChatParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostCloseChat(ctx context.Context, params *PostCloseChatParams, body PostCloseChatJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostFreeHands request
	PostFreeHands(ctx context.Context, params *PostFreeHandsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostGetChatHistoryWithBody request with any body
	PostGetChatHistoryWithBody(ctx context.Context, params *PostGetChatHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostGetChatHistory(ctx context.Context, params *PostGetChatHistoryParams, body PostGetChatHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostGetChats request
	PostGetChats(ctx context.Context, params *PostGetChatsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostGetFreeHandsBtnAvailability request
	PostGetFreeHandsBtnAvailability(ctx context.Context, params *PostGetFreeHandsBtnAvailabilityParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostSendMessageWithBody request with any body
	PostSendMessageWithBody(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostSendMessage(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostCloseChatWithBody(ctx context.Context, params *PostCloseChatParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostCloseChatRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostCloseChat(ctx context.Context, params *PostCloseChatParams, body PostCloseChatJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostCloseChatRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostFreeHands(ctx context.Context, params *PostFreeHandsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostFreeHandsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostGetChatHistoryWithBody(ctx context.Context, params *PostGetChatHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostGetChatHistoryRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostGetChatHistory(ctx context.Context, params *PostGetChatHistoryParams, body PostGetChatHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostGetChatHistoryRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostGetChats(ctx context.Context, params *PostGetChatsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostGetChatsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostGetFreeHandsBtnAvailability(ctx context.Context, params *PostGetFreeHandsBtnAvailabilityParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostGetFreeHandsBtnAvailabilityRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostSendMessageWithBody(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostSendMessageRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostSendMessage(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostSendMessageRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostCloseChatRequest calls the generic PostCloseChat builder with application/json body
func NewPostCloseChatRequest(server string, params *PostCloseChatParams, body PostCloseChatJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostCloseChatRequestWithBody(server, params, "application/json", bodyReader)
}

// NewPostCloseChatRequestWithBody generates requests for PostCloseChat with any type of body
func NewPostCloseChatRequestWithBody(server string, params *PostCloseChatParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/closeChat")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

// NewPostFreeHandsRequest generates requests for PostFreeHands
func NewPostFreeHandsRequest(server string, params *PostFreeHandsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/freeHands")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

// NewPostGetChatHistoryRequest calls the generic PostGetChatHistory builder with application/json body
func NewPostGetChatHistoryRequest(server string, params *PostGetChatHistoryParams, body PostGetChatHistoryJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostGetChatHistoryRequestWithBody(server, params, "application/json", bodyReader)
}

// NewPostGetChatHistoryRequestWithBody generates requests for PostGetChatHistory with any type of body
func NewPostGetChatHistoryRequestWithBody(server string, params *PostGetChatHistoryParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/getChatHistory")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

// NewPostGetChatsRequest generates requests for PostGetChats
func NewPostGetChatsRequest(server string, params *PostGetChatsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/getChats")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

// NewPostGetFreeHandsBtnAvailabilityRequest generates requests for PostGetFreeHandsBtnAvailability
func NewPostGetFreeHandsBtnAvailabilityRequest(server string, params *PostGetFreeHandsBtnAvailabilityParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/getFreeHandsBtnAvailability")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

// NewPostSendMessageRequest calls the generic PostSendMessage builder with application/json body
func NewPostSendMessageRequest(server string, params *PostSendMessageParams, body PostSendMessageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostSendMessageRequestWithBody(server, params, "application/json", bodyReader)
}

// NewPostSendMessageRequestWithBody generates requests for PostSendMessage with any type of body
func NewPostSendMessageRequestWithBody(server string, params *PostSendMessageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/sendMessage")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PostCloseChatWithBodyWithResponse request with any body
	PostCloseChatWithBodyWithResponse(ctx context.Context, params *PostCloseChatParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCloseChatResponse, error)

	PostCloseChatWithResponse(ctx context.Context, params *PostCloseChatParams, body PostCloseChatJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCloseChatResponse, error)

	// PostFreeHandsWithResponse request
	PostFreeHandsWithResponse(ctx context.Context, params *PostFreeHandsParams, reqEditors ...RequestEditorFn) (*PostFreeHandsResponse, error)

	// PostGetChatHistoryWithBodyWithResponse request with any body
	PostGetChatHistoryWithBodyWithResponse(ctx context.Context, params *PostGetChatHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostGetChatHistoryResponse, error)

	PostGetChatHistoryWithResponse(ctx context.Context, params *PostGetChatHistoryParams, body PostGetChatHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*PostGetChatHistoryResponse, error)

	// PostGetChatsWithResponse request
	PostGetChatsWithResponse(ctx context.Context, params *PostGetChatsParams, reqEditors ...RequestEditorFn) (*PostGetChatsResponse, error)

	// PostGetFreeHandsBtnAvailabilityWithResponse request
	PostGetFreeHandsBtnAvailabilityWithResponse(ctx context.Context, params *PostGetFreeHandsBtnAvailabilityParams, reqEditors ...RequestEditorFn) (*PostGetFreeHandsBtnAvailabilityResponse, error)

	// PostSendMessageWithBodyWithResponse request with any body
	PostSendMessageWithBodyWithResponse(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error)

	PostSendMessageWithResponse(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error)
}

type PostCloseChatResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *CloseChatResponse
}

// Status returns HTTPResponse.Status
func (r PostCloseChatResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostCloseChatResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostFreeHandsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *FreeHandsResponse
}

// Status returns HTTPResponse.Status
func (r PostFreeHandsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostFreeHandsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostGetChatHistoryResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetChatHistoryResponse
}

// Status returns HTTPResponse.Status
func (r PostGetChatHistoryResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostGetChatHistoryResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostGetChatsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetChatsResponse
}

// Status returns HTTPResponse.Status
func (r PostGetChatsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostGetChatsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostGetFreeHandsBtnAvailabilityResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetFreeHandsBtnAvailabilityResponse
}

// Status returns HTTPResponse.Status
func (r PostGetFreeHandsBtnAvailabilityResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostGetFreeHandsBtnAvailabilityResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostSendMessageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SendMessageResponse
}

// Status returns HTTPResponse.Status
func (r PostSendMessageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostSendMessageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostCloseChatWithBodyWithResponse request with arbitrary body returning *PostCloseChatResponse
func (c *ClientWithResponses) PostCloseChatWithBodyWithResponse(ctx context.Context, params *PostCloseChatParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCloseChatResponse, error) {
	rsp, err := c.PostCloseChatWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostCloseChatResponse(rsp)
}

func (c *ClientWithResponses) PostCloseChatWithResponse(ctx context.Context, params *PostCloseChatParams, body PostCloseChatJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCloseChatResponse, error) {
	rsp, err := c.PostCloseChat(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostCloseChatResponse(rsp)
}

// PostFreeHandsWithResponse request returning *PostFreeHandsResponse
func (c *ClientWithResponses) PostFreeHandsWithResponse(ctx context.Context, params *PostFreeHandsParams, reqEditors ...RequestEditorFn) (*PostFreeHandsResponse, error) {
	rsp, err := c.PostFreeHands(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostFreeHandsResponse(rsp)
}

// PostGetChatHistoryWithBodyWithResponse request with arbitrary body returning *PostGetChatHistoryResponse
func (c *ClientWithResponses) PostGetChatHistoryWithBodyWithResponse(ctx context.Context, params *PostGetChatHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostGetChatHistoryResponse, error) {
	rsp, err := c.PostGetChatHistoryWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostGetChatHistoryResponse(rsp)
}

func (c *ClientWithResponses) PostGetChatHistoryWithResponse(ctx context.Context, params *PostGetChatHistoryParams, body PostGetChatHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*PostGetChatHistoryResponse, error) {
	rsp, err := c.PostGetChatHistory(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostGetChatHistoryResponse(rsp)
}

// PostGetChatsWithResponse request returning *PostGetChatsResponse
func (c *ClientWithResponses) PostGetChatsWithResponse(ctx context.Context, params *PostGetChatsParams, reqEditors ...RequestEditorFn) (*PostGetChatsResponse, error) {
	rsp, err := c.PostGetChats(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostGetChatsResponse(rsp)
}

// PostGetFreeHandsBtnAvailabilityWithResponse request returning *PostGetFreeHandsBtnAvailabilityResponse
func (c *ClientWithResponses) PostGetFreeHandsBtnAvailabilityWithResponse(ctx context.Context, params *PostGetFreeHandsBtnAvailabilityParams, reqEditors ...RequestEditorFn) (*PostGetFreeHandsBtnAvailabilityResponse, error) {
	rsp, err := c.PostGetFreeHandsBtnAvailability(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostGetFreeHandsBtnAvailabilityResponse(rsp)
}

// PostSendMessageWithBodyWithResponse request with arbitrary body returning *PostSendMessageResponse
func (c *ClientWithResponses) PostSendMessageWithBodyWithResponse(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error) {
	rsp, err := c.PostSendMessageWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostSendMessageResponse(rsp)
}

func (c *ClientWithResponses) PostSendMessageWithResponse(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error) {
	rsp, err := c.PostSendMessage(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostSendMessageResponse(rsp)
}

// ParsePostCloseChatResponse parses an HTTP response from a PostCloseChatWithResponse call
func ParsePostCloseChatResponse(rsp *http.Response) (*PostCloseChatResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostCloseChatResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest CloseChatResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostFreeHandsResponse parses an HTTP response from a PostFreeHandsWithResponse call
func ParsePostFreeHandsResponse(rsp *http.Response) (*PostFreeHandsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostFreeHandsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest FreeHandsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostGetChatHistoryResponse parses an HTTP response from a PostGetChatHistoryWithResponse call
func ParsePostGetChatHistoryResponse(rsp *http.Response) (*PostGetChatHistoryResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostGetChatHistoryResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetChatHistoryResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostGetChatsResponse parses an HTTP response from a PostGetChatsWithResponse call
func ParsePostGetChatsResponse(rsp *http.Response) (*PostGetChatsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostGetChatsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetChatsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostGetFreeHandsBtnAvailabilityResponse parses an HTTP response from a PostGetFreeHandsBtnAvailabilityWithResponse call
func ParsePostGetFreeHandsBtnAvailabilityResponse(rsp *http.Response) (*PostGetFreeHandsBtnAvailabilityResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostGetFreeHandsBtnAvailabilityResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetFreeHandsBtnAvailabilityResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostSendMessageResponse parses an HTTP response from a PostSendMessageWithResponse call
func ParsePostSendMessageResponse(rsp *http.Response) (*PostSendMessageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostSendMessageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SendMessageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xYXW/bOhL9K8TsPuwCsiU3W6AwsA9p0za5aNqg6UUL5PqBlsYWG4pUyZGb3ED//YLU",
	"h2VLjtO0KdI+JZaG5MyZM8MzuoFYZ7lWqMjC9AZybniGhMb/+vQevxRo6eToGHmCxj0TCqaQVj8DUDxD",
	"mMKnUW05OjmCAAx+KYTBBKZkCgzAxilm3K1eaJNxgikUhUggALrO3XpLRqglBHA1WupR/dD9sePWhe7b",
	"kchybajymFKYwlJQWszHsc5CvExXXOpVOOfqchSnnEYWzUrEGApFaBSXod8byrIsG+98wC9S7jflUr5b",
	"wPTiBv5tcAFT+Fe4ximsF4TO+iSBMriB3OgcDQn028RSoHKv7hXxnxbNQ4Xbzc3F2s9Z65eef8aYoJyV",
	"AdTxTXvhtc+/PTi/588JrvKyCeSNsDQciv9HEGb+n30Jh7INkxvDr2HoXFsdK7VFt6bm8O+A5Dokm2tl",
	"sR9TwsmVehkAGqPNPkxfeiN/5svGfgskneCddnnhDMsAEiQupF+7CWMZQIbW8iUOvNsKujEMqvNnjX8v",
	"am8StLEROQntGmKsFXGhLDv+8OGM+cCZW2cZVwmzOcZiIWI2L6xQaC2TeiniDbv/UIpMckssKyyxObK/",
	"iig6wP+zSRRF/x1DAKiKDKYXT6MoCp5G0WQWQCaUyNzT/0VRSxqXz6Vvz1cjt2a04sY1auviaoM45Yov",
	"0bxboZGaJ5hU8HfCfKsPYxIrPDN6LjE7UZ7/DolXBvGYq8QerriQfC6koOt+5nj1Vnbhnmstkase3mvb",
	"jQMegmevkVwkx8KSNte/WnEGEBfGVgH36J3zJZ6Lvz1cGb+quDFx3GiZMukT5ZaC38ZqTzZuTcFpVVL2",
	"zNXV/fNmv8+L9jK4nwctM5+T6rL/+5waLqh7eHi67m9btVhQqs0j1CQBzHVyPcjm2CAnTA5pw+mEE45I",
	"ZNjzvAxA3DPAGrcHi1HY9xijWGFyh17ofW4TVgPUhWNjv9k67R8Fpbqg5zWgvwwDfptE70nkOs5O0qp+",
	"2EtXLUDurk6b0u8J1AAUXtGdJY+FeoHz8W0h5REn/hD38DmqpHa6cwl/5+xVB9GUQMav3qBaupweRPU9",
	"2DyYBHcDxO81PCJthPAD7sZuAX8zom6Wxbgwgq7P3bvq9Dlyg+awcBE3v141ZfPHxw9QT8C+Ifm36zpK",
	"ifIqV0IttGeQICfn4DlXl+y8yF3ZMJcMVutJdnh2AgGs0NhKGq8mLhKdo+K5gCkcjKPxAQS+0LyDYdwM",
	"FR46XdFgU1+fcnPJ8kqHMm6Z1XKFiRfXfjVzZQn+GMPdGtfm4ExbaicWf+T668YOfq1Nwt7Xj3JW0QNt",
	"22Gd7kdV8TbPpYj94eFn67y+6Xz4uJXL23PiVh8hU6B/UBHMg/Ykih7i/JrC3oHNDPgce6yTmmvhohEt",
	"u/N2hLHkBplBnlSTD2mWcpVIbNJpx4N5axXRj8rbA4HXn1QGwKtg8Nsz646sEVxuqOvdML5G8gxnaWU5",
	"DNmmVn+8fB+ev34y6XcMNgPJa65pJoWl8Vbq7O1J80O9sMT0wifQsq+CUua64R7+N8POI6d/bybb0Tr6",
	"6O0apG4HlHcsO71E4df9eO488dFDvHfoHEC9a8cM2kK2bceuVctutJ208bjWKsih7ejsaDwMcUcMPd7m",
	"MyA6f3LnGdKMu9sOq2eGqnQ6Es+j2hV3FzOHmRtNGsy3L6EVSp1nqIhVVhBAYWSt86ZhKHXMZaotTZ9F",
	"zyahV2637YJPcHsnv2hW/hMAAP//co4OD1gaAAA=",
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
