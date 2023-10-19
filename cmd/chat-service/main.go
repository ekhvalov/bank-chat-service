package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ekhvalov/bank-chat-service/internal/buildinfo"
	keycloakclient "github.com/ekhvalov/bank-chat-service/internal/clients/keycloak"
	"github.com/ekhvalov/bank-chat-service/internal/config"
	"github.com/ekhvalov/bank-chat-service/internal/logger"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
	serverdebug "github.com/ekhvalov/bank-chat-service/internal/server-debug"
	"github.com/ekhvalov/bank-chat-service/internal/store"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	mustInitGlobalLogger(cfg)
	defer logger.Sync()

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %v", err)
	}
	keycloakClient, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode),
		keycloakclient.WithUserAgent(fmt.Sprintf("chat-service/%s", buildinfo.BuildInfo.Main.Version)),
	))
	if err != nil {
		return fmt.Errorf("keycloak client create: %v", err)
	}
	if cfg.Global.IsProduction() && cfg.Clients.Keycloak.DebugMode {
		zap.L().Warn("keycloak client debug mode enabled on production environment")
	}
	if cfg.Global.IsProduction() && cfg.Clients.Postgres.DebugMode {
		zap.L().Warn("postgres client debug mode enabled on production environment")
	}

	msgRepo := mustInitMsgRepo(ctx, cfg.Clients.Postgres)

	srvClient, err := initServerClient(
		cfg.Servers.Client,
		swagger,
		keycloakClient,
		msgRepo,
		cfg.Global.IsProduction(),
	)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return srvClient.Run(ctx) })

	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}

func mustInitGlobalLogger(cfg config.Config) {
	logger.MustInit(logger.NewOptions(
		cfg.Log.Level,
		cfg.Global.Env,
		logger.WithProductionMode(cfg.Global.IsProduction()),
		logger.WithSentryDSN(cfg.Sentry.DSN),
	))
}

func mustInitMsgRepo(ctx context.Context, cfg config.PostgresClientConfig) *messagesrepo.Repo {
	storeClient, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.Address,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		store.WithDebug(cfg.DebugMode),
	))
	if err != nil {
		panic(fmt.Sprintf("create psql store client: %v", err))
	}
	if err = storeClient.Schema.Create(ctx); err != nil {
		panic(fmt.Sprintf("create schema: %v", err))
	}

	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(store.NewDatabase(storeClient)))
	if err != nil {
		panic(fmt.Sprintf("create messages repo: %v", err))
	}

	return msgRepo
}
