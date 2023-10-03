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
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ekhvalov/bank-chat-service/internal/buildinfo"
	"github.com/ekhvalov/bank-chat-service/internal/logger"
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
	e.Use(middleware.Recover())

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

	e.PUT("/log/level", s.logLevel)

	s.debugPprof(e.Group("/debug/pprof"))

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
	g.GET("/heap", func(ctx echo.Context) error {
		pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})
	g.GET("/goroutine", func(ctx echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	})
	g.GET("/block", func(ctx echo.Context) error {
		pprof.Handler("block").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	})
	g.GET("/threadcreate", func(ctx echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	})
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
	g.GET("/mutex", func(ctx echo.Context) error {
		pprof.Handler("mutex").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	})
}
