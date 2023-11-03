package logger

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ekhvalov/bank-chat-service/internal/buildinfo"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var Levels = []string{LevelDebug, LevelInfo, LevelWarn, LevelError}

var atomicLevel zap.AtomicLevel

type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	env            string `option:"mandatory" validate:"required,oneof=dev stage prod"`
	sentryDSN      string `validate:"omitempty,http_url"`
	productionMode bool
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate options: %v", err)
	}
	cfg := zap.NewProductionEncoderConfig()
	cfg.NameKey = "component"
	cfg.TimeKey = "T"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	var encoder zapcore.Encoder
	if opts.productionMode {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(cfg)
	} else {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(cfg)
	}

	lvl, err := zapcore.ParseLevel(opts.level)
	if err != nil {
		return fmt.Errorf("log level error: %v", err)
	}
	atomicLevel = zap.NewAtomicLevelAt(lvl)
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, atomicLevel),
	}
	if opts.sentryDSN != "" {
		sCfg := zapsentry.Configuration{
			Level: zapcore.WarnLevel,
		}
		client, errSentry := NewSentryClient(opts.sentryDSN, opts.env, buildinfo.BuildInfo.Main.Version)
		if errSentry != nil {
			return fmt.Errorf("sentry client: %v", errSentry)
		}
		core, errCore := zapsentry.NewCore(sCfg, zapsentry.NewSentryClientFromClient(client))
		if errCore != nil {
			return fmt.Errorf("zapsentry create core: %v", errCore)
		}
		cores = append(cores, core)
	}
	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !(errors.Is(err, syscall.ENOTTY) || errors.Is(err, syscall.EINVAL)) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}

func ChangeLevel(level string) error {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		return err
	}
	atomicLevel.SetLevel(lvl)
	return nil
}
