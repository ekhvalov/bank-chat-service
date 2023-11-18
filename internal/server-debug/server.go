package serverdebug

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ekhvalov/bank-chat-service/internal/buildinfo"
	"github.com/ekhvalov/bank-chat-service/internal/logger"
	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	clientevents "github.com/ekhvalov/bank-chat-service/internal/server-client/events"
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

type Options struct {
	addr string `option:"mandatory" validate:"required,hostname_port"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options error: %v", err)
	}

	lg := zap.L().Named("server-debug")

	e := echo.New()
	e.Use(
		middlewares.NewRecover(lg),
		middlewares.NewLogger(lg),
	)

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
	index := newIndexPage()

	e.GET("/version", s.version)
	index.addPage("/version", "Get build information")
	index.addPage("/debug/pprof/", "Go std profiler")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half-min profile")
	index.addPage("/debug/error", "Test log error")
	index.addPage("/schema/client", "Get client OpenAPI specification")
	index.addPage("/schema/client-events", "Get OpenAPI specification of events for client")
	index.addPage("/schema/manager", "Get manager OpenAPI specification")

	e.PUT("/log/level", s.logLevel)

	s.debugPprof(e.Group("/debug/pprof"))
	e.GET("/debug/error", s.debugError)
	e.GET("/schema/client", s.schemaClient)
	e.GET("/schema/client-events", s.schemaClientEvents)
	e.GET("/schema/manager", s.schemaManager)

	e.GET("/", index.handler)
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

func (s *Server) version(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, buildinfo.BuildInfo)
}

func (s *Server) logLevel(eCtx echo.Context) error {
	level := eCtx.FormValue("level")
	if len(level) == 0 {
		return eCtx.String(http.StatusBadRequest, "log level is empty")
	}
	if err := logger.ChangeLevel(level); err != nil {
		return eCtx.String(http.StatusBadRequest, fmt.Sprintf("change log level error: %q", level))
	}
	return eCtx.String(http.StatusOK, "OK")
}

func (s *Server) debugPprof(g *echo.Group) {
	indexHandler := func(ctx echo.Context) error {
		pprof.Index(ctx.Response().Writer, ctx.Request())
		return nil
	}
	g.GET("", indexHandler)
	g.GET("/", indexHandler)
	g.GET("/heap", echo.WrapHandler(pprof.Handler("heap")))
	g.GET("/goroutine", echo.WrapHandler(pprof.Handler("goroutine")))
	g.GET("/block", echo.WrapHandler(pprof.Handler("block")))
	g.GET("/threadcreate", echo.WrapHandler(pprof.Handler("threadcreate")))
	g.GET("/cmdline", func(ctx echo.Context) error {
		pprof.Cmdline(ctx.Response().Writer, ctx.Request())
		return nil
	})
	g.GET("/profile", func(ctx echo.Context) error {
		pprof.Profile(ctx.Response().Writer, ctx.Request())
		return nil
	})
	symbolHandler := func(ctx echo.Context) error {
		pprof.Symbol(ctx.Response().Writer, ctx.Request())
		return nil
	}
	g.GET("/symbol", symbolHandler)
	g.POST("/symbol", symbolHandler)
	g.GET("/trace", func(ctx echo.Context) error {
		pprof.Trace(ctx.Response().Writer, ctx.Request())
		return nil
	})
	g.GET("/mutex", echo.WrapHandler(pprof.Handler("mutex")))
}

func (s *Server) debugError(eCtx echo.Context) error {
	s.lg.Error("test error")
	return eCtx.String(http.StatusOK, "OK")
}

func (s *Server) schemaClient(eCtx echo.Context) error {
	sw, err := clientv1.GetSwagger()
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, fmt.Sprintf("get swagger error: %v", err))
	}
	return eCtx.JSON(http.StatusOK, sw)
}

func (s *Server) schemaClientEvents(eCtx echo.Context) error {
	sw, err := clientevents.GetSwagger()
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, fmt.Sprintf("get swagger error: %v", err))
	}
	return eCtx.JSON(http.StatusOK, sw)
}

func (s *Server) schemaManager(eCtx echo.Context) error {
	sw, err := managerv1.GetSwagger()
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, fmt.Sprintf("get swagger error: %v", err))
	}
	return eCtx.JSON(http.StatusOK, sw)
}
