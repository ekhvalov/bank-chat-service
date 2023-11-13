package serverclient

import (
	"fmt"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"

	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	v1Swagger  *openapi3.T       `option:"mandatory" validate:"required"`
	v1Handlers clientv1.Handlers `option:"mandatory" validate:"required"`
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Server{Options: opts}, nil
}

type Server struct {
	Options
}

func (r *Server) RegisterHandlers(e *echo.Echo) {
	v1 := e.Group("v1", oapimdlwr.OapiRequestValidatorWithOptions(r.v1Swagger, &oapimdlwr.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody:  false,
			ExcludeResponseBody: true,
			AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
		},
	}))
	clientv1.RegisterHandlers(v1, r.v1Handlers)
}
