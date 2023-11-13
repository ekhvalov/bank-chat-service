package websocketstream

import (
	"context"
	"fmt"
	"time"

	gorillaws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const (
	writeTimeout = time.Second
)

type eventStream interface {
	Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error)
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	pingPeriod time.Duration `default:"3s" validate:"omitempty,min=100ms,max=30s"`

	logger       *zap.Logger     `option:"mandatory" validate:"required"`
	eventStream  eventStream     `option:"mandatory" validate:"required"`
	eventAdapter EventAdapter    `option:"mandatory" validate:"required"`
	eventWriter  EventWriter     `option:"mandatory" validate:"required"`
	upgrader     Upgrader        `option:"mandatory" validate:"required"`
	shutdownCh   <-chan struct{} `option:"mandatory" validate:"required"`
}

type HTTPHandler struct {
	Options
}

func NewHTTPHandler(opts Options) (*HTTPHandler, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &HTTPHandler{Options: opts}, nil
}

func (h *HTTPHandler) Serve(eCtx echo.Context) error {
	ws, err := h.upgrader.Upgrade(eCtx.Response(), eCtx.Request(), nil)
	if err != nil {
		return fmt.Errorf("upgrade to websocket: %v", err)
	}
	uid := middlewares.MustUserID(eCtx)
	events, err := h.eventStream.Subscribe(eCtx.Request().Context(), uid)
	if err != nil {
		return fmt.Errorf("events subscribe: %v", err)
	}
	eg := &errgroup.Group{}

	eg.Go(func() error {
		return h.readLoop(eCtx.Request().Context(), ws)
	})

	eg.Go(func() error {
		return h.writeLoop(eCtx.Request().Context(), ws, events)
	})

	eg.Go(func() error {
		closer := newWsCloser(h.logger, ws)
		<-h.shutdownCh
		closer.Close(gorillaws.CloseNormalClosure)
		return nil
	})

	return eg.Wait()
}

// readLoop listen PONGs.
func (h *HTTPHandler) readLoop(_ context.Context, ws Websocket) error {
	for {
		ws.SetPongHandler(nil)
		msgType, r, err := ws.NextReader()
		if err != nil {
			return fmt.Errorf("get message reader: %v", err)
		}
		if msgType == gorillaws.PongMessage {
			msg := make([]byte, 0)
			_, err := r.Read(msg)
			if err != nil {
				return fmt.Errorf("read message: %v", err)
			}
			h.logger.Debug(string(msg))
		}
	}
}

// writeLoop listen events and writes them into Websocket.
func (h *HTTPHandler) writeLoop(ctx context.Context, ws Websocket, events <-chan eventstream.Event) error {
	jsonWriter := JSONEventWriter{}
	ticker := time.NewTicker(h.pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := ws.WriteControl(gorillaws.PingMessage, nil, time.Now().Add(writeTimeout))
			if err != nil {
				return fmt.Errorf("write ping messa: %v", err)
			}
		case e, ok := <-events:
			if !ok {
				return nil
			}
			w, err := ws.NextWriter(gorillaws.TextMessage)
			if err != nil {
				return fmt.Errorf("get message writer: %v", err)
			}
			err = jsonWriter.Write(e, w)
			if errClose := w.Close(); errClose != nil {
				h.logger.Debug("close writer", zap.Error(errClose))
			}
			if err != nil {
				return fmt.Errorf("event write: %v", err)
			}
		}
	}
}
