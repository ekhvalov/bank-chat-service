package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	serverclient "github.com/ekhvalov/bank-chat-service/internal/server-client"
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
	servermanager "github.com/ekhvalov/bank-chat-service/internal/server-manager"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
	websocketstream "github.com/ekhvalov/bank-chat-service/internal/websocket-stream"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
	// bodyLimit 3000 unicode symbols of 4 bytes each.
	bodyLimit = "12KB"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr           string                       `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins   []string                     `option:"mandatory" validate:"min=1"`
	accessResource string                       `option:"mandatory" validate:"required"`
	accessRole     string                       `option:"mandatory" validate:"required"`
	secWsProtocol  string                       `option:"mandatory" validate:"required"`
	logger         *zap.Logger                  `option:"mandatory" validate:"required"`
	jwtParser      *middlewares.JWTParser       `option:"mandatory" validate:"required"`
	errorHandler   echo.HTTPErrorHandler        `option:"mandatory" validate:"required"`
	wsHandler      *websocketstream.HTTPHandler `option:"mandatory" validate:"required"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

type Handlers interface {
	clientv1.Handlers | managerv1.Handlers
	RegisterTo(e *echo.Echo) error
}

type HandlersRegistrar interface {
	*servermanager.Server | *serverclient.Server
	RegisterHandlers(e *echo.Echo)
}

func New[T HandlersRegistrar](opts Options, registrar T) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options error: %v", err)
	}
	if nil == registrar {
		return nil, fmt.Errorf("expected registrar of type %T, got <nil>", registrar)
	}

	e := echo.New()
	e.HTTPErrorHandler = opts.errorHandler
	e.Use(
		middlewares.NewRecover(opts.logger),
		middlewares.NewLogger(opts.logger),
		middleware.BodyLimit(bodyLimit),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: opts.allowOrigins,
			AllowMethods: []string{http.MethodPost},
		}),
		middlewares.NewKeycloakTokenAuth(
			opts.jwtParser,
			opts.accessResource,
			opts.accessRole,
			opts.secWsProtocol,
		),
	)
	registrar.RegisterHandlers(e)

	e.GET("/ws", opts.wsHandler.Serve)

	s := &Server{
		lg: opts.logger,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}
