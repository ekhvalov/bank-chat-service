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

	"github.com/ekhvalov/bank-chat-service/internal/config"
	"github.com/ekhvalov/bank-chat-service/internal/logger"
	jobsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/jobs"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	clientevents "github.com/ekhvalov/bank-chat-service/internal/server-client/events"
	serverdebug "github.com/ekhvalov/bank-chat-service/internal/server-debug"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	inmemeventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream/in-mem"
	msgproducer "github.com/ekhvalov/bank-chat-service/internal/services/msg-producer"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	sendclientmessagejob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-client-message"
	"github.com/ekhvalov/bank-chat-service/internal/store"
	storegen "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	websocketstream "github.com/ekhvalov/bank-chat-service/internal/websocket-stream"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

const (
	logNameMain        = "main"
	logNameMsgProducer = "msg-producer"
	logNameOutbox      = "outbox"
)

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

	lg := zap.L().Named(logNameMain)
	if cfg.Global.IsProduction() && cfg.Clients.Keycloak.DebugMode {
		lg.Warn("keycloak client in the debug mode")
	}
	if cfg.Global.IsProduction() && cfg.Clients.Postgres.DebugMode {
		lg.Warn("psql client in the debug mode")
	}
	storeClient, err := initStoreClient(ctx, cfg.Clients.Postgres)
	if err != nil {
		return fmt.Errorf("init store client: %v", err)
	}
	defer func() {
		for _, shutdownFunc := range shutdownFunctions {
			if err := shutdownFunc(); err != nil {
				errReturned = errors.Join(errReturned, err)
			}
		}
	}()

	storeDB := storegen.NewDatabase(storeClient)

	msgProducerSvc, err := initMsgProducerService(cfg.Services.MsgProducer)
	if err != nil {
		return fmt.Errorf("inig msg producer service: %v", err)
	}

	eventStream := inmemeventstream.New()

	outboxSvc, err := initOutboxService(cfg.Services.OutboxService, storeDB, msgProducerSvc, eventStream)
	if err != nil {
		return fmt.Errorf("init outbox service: %v", err)
	}

	srvClient, err := initServerClient(cfg, storeDB, outboxSvc, eventStream)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	srvManager, err := initServerManager(cfg, storeDB, eventStream)
	if err != nil {
		return fmt.Errorf("init manager server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return srvClient.Run(ctx) })
	eg.Go(func() error { return srvManager.Run(ctx) })

	// Run services.
	eg.Go(func() error { return outboxSvc.Run(ctx) })

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

func initStoreClient(ctx context.Context, cfg config.PostgresClientConfig) (*storegen.Client, error) {
	client, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.Address,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		store.WithDebug(cfg.DebugMode),
	))
	if err != nil {
		return nil, fmt.Errorf("create psql client: %v", err)
	}

	if err = client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("create schema: %v", err)
	}

	registerShutdownFunc(func() error {
		if err := client.Close(); err != nil {
			return fmt.Errorf("close pgsql: %v", err)
		}
		return nil
	})

	return client, nil
}

func initOutboxService(
	cfg config.OutboxService,
	storeDB *storegen.Database,
	msgProducer *msgproducer.Service,
	eventStream eventstream.EventStream,
) (*outbox.Service, error) {
	repo, err := jobsrepo.New(jobsrepo.NewOptions(storeDB))
	if err != nil {
		return nil, fmt.Errorf("create jobs repo: %v", err)
	}
	lg := zap.L().Named(logNameOutbox)
	outboxSvc, err := outbox.New(outbox.NewOptions(cfg.Workers, cfg.IdleTime, cfg.ReserveFor, repo, storeDB, lg))
	if err != nil {
		return nil, fmt.Errorf("create outbox service: %v", err)
	}
	messagesRepo, err := messagesrepo.New(messagesrepo.NewOptions(storeDB))
	if err != nil {
		return nil, fmt.Errorf("create messages repo: %v", err)
	}

	sendClientMsgJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(
		msgProducer,
		messagesRepo,
		eventStream,
		lg,
	))
	if err != nil {
		return nil, fmt.Errorf("create send-client-message job: %v", err)
	}
	err = outboxSvc.RegisterJob(sendClientMsgJob)
	if err != nil {
		return nil, fmt.Errorf("register send-client-message job: %v", err)
	}

	return outboxSvc, nil
}

func initMsgProducerService(cfg config.MsgProducerServiceConfig) (*msgproducer.Service, error) {
	kw := msgproducer.NewKafkaWriter(cfg.Brokers, cfg.Topic, cfg.BatchSize)
	lg := zap.L().Named(logNameMsgProducer)
	msgProducer, err := msgproducer.New(msgproducer.NewOptions(kw, lg, msgproducer.WithEncryptKey(cfg.EncryptKey)))
	if err != nil {
		return nil, fmt.Errorf("create msg producer: %v", err)
	}
	return msgProducer, nil
}

func initWebsocketHandler(lg *zap.Logger, allowOrigins []string, secWsProtocol string, eventStream eventstream.EventStream) (
	*websocketstream.HTTPHandler,
	error,
) {
	shutdownCh := make(chan struct{})
	wsHandler, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		lg,
		eventStream,
		clientevents.Adapter{},
		websocketstream.JSONEventWriter{},
		websocketstream.NewUpgrader(allowOrigins, secWsProtocol),
		shutdownCh,
	))
	registerShutdownFunc(func() error {
		close(shutdownCh)
		return nil
	})
	return wsHandler, err
}
