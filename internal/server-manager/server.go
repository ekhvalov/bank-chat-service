package servermanager

import (
	"context"
	"fmt"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	"github.com/ekhvalov/bank-chat-service/internal/server"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr           string                    `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins   []string                  `option:"mandatory" validate:"min=1"`
	accessResource string                    `option:"mandatory" validate:"required"`
	accessRole     string                    `option:"mandatory" validate:"required"`
	logger         *zap.Logger               `option:"mandatory" validate:"required"`
	v1Swagger      *openapi3.T               `option:"mandatory" validate:"required"`
	v1Handlers     managerv1.ServerInterface `option:"mandatory" validate:"required"`
	introspector   middlewares.Introspector  `option:"mandatory" validate:"required"`
	errorHandler   echo.HTTPErrorHandler     `option:"mandatory" validate:"required"`
}

type Server struct {
	srv *server.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options error: %v", err)
	}

	handlersRegistrar := func(e *echo.Echo) {
		v1 := e.Group("v1", oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
			Options: openapi3filter.Options{
				ExcludeRequestBody:  false,
				ExcludeResponseBody: true,
				AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
			},
		}))
		managerv1.RegisterHandlers(v1, opts.v1Handlers)
	}
	srv, err := server.New(server.NewOptions(
		opts.addr,
		opts.allowOrigins,
		opts.accessResource,
		opts.accessRole,
		opts.logger,
		opts.v1Swagger,
		opts.introspector,
		opts.errorHandler,
		handlersRegistrar,
	))
	if err != nil {
		return nil, fmt.Errorf("create server: %v", err)
	}

	return &Server{srv: srv}, nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.srv.Run(ctx)
}
