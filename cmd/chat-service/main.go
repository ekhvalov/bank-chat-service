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
	chatsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/chats"
	jobsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/jobs"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	serverdebug "github.com/ekhvalov/bank-chat-service/internal/server-debug"
	afcverdictsprocessor "github.com/ekhvalov/bank-chat-service/internal/services/afc-verdicts-processor"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	inmemeventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream/in-mem"
	managerload "github.com/ekhvalov/bank-chat-service/internal/services/manager-load"
	managerpool "github.com/ekhvalov/bank-chat-service/internal/services/manager-pool"
	inmemmanagerpool "github.com/ekhvalov/bank-chat-service/internal/services/manager-pool/in-mem"
	managerscheduler "github.com/ekhvalov/bank-chat-service/internal/services/manager-scheduler"
	msgproducer "github.com/ekhvalov/bank-chat-service/internal/services/msg-producer"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	clientmessageblockedjob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/client-message-sent"
	closechatjob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/close-chat"
	managerassignedtoproblemjob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	sendclientmessagejob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-client-message"
	sendmanagermessage "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/ekhvalov/bank-chat-service/internal/store"
	storegen "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	websocketstream "github.com/ekhvalov/bank-chat-service/internal/websocket-stream"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

const (
	logNameMain         = "main"
	logNameMsgProducer  = "msg-producer"
	logNameOutbox       = "outbox"
	logNameAFCProcessor = "afc-processor"
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

	messagesRepo, err := messagesrepo.New(messagesrepo.NewOptions(storeDB))
	if err != nil {
		return fmt.Errorf("create messages repo: %v", err)
	}

	chatsRepo, err := chatsrepo.New(chatsrepo.NewOptions(storeDB))
	if err != nil {
		return fmt.Errorf("create chats repo: %v", err)
	}

	problemsRepo, err := problemsrepo.New(problemsrepo.NewOptions(storeDB))
	if err != nil {
		return fmt.Errorf("create problems repo: %v", err)
	}

	managerLoad, err := managerload.New(managerload.NewOptions(cfg.Services.ManagerLoad.MaxProblemsAtSameTime, problemsRepo))
	if err != nil {
		return fmt.Errorf("create manager load service: %v", err)
	}

	outboxSvc, err := initOutboxService(
		cfg.Services.OutboxService,
		storeDB,
		msgProducerSvc,
		messagesRepo,
		chatsRepo,
		problemsRepo,
		managerLoad,
		eventStream,
	)
	if err != nil {
		return fmt.Errorf("init outbox service: %v", err)
	}

	serverClient, err := initServerClient(cfg, storeDB, outboxSvc, eventStream)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	managerPool := inmemmanagerpool.New()

	serverManager, err := initServerManager(
		cfg,
		eventStream,
		managerPool,
		chatsRepo,
		messagesRepo,
		problemsRepo,
		outboxSvc,
		storeDB,
	)
	if err != nil {
		return fmt.Errorf("init manager server: %v", err)
	}

	afcProcessor, err := initAFCVerdictsProcessor(cfg.Services.AFCVerdictsProcessorService, storeDB, messagesRepo, outboxSvc)
	if err != nil {
		return fmt.Errorf("init AFC verdicts processor: %v", err)
	}

	managerScheduler, err := initManagerScheduler(
		cfg.Services.ManagerSchedulerService,
		managerPool,
		messagesRepo,
		outboxSvc,
		problemsRepo,
		storeDB,
		lg,
	)
	if err != nil {
		return fmt.Errorf("init manager scheduler: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return serverClient.Run(ctx) })
	eg.Go(func() error { return serverManager.Run(ctx) })

	// Run services.
	eg.Go(func() error { return outboxSvc.Run(ctx) })
	eg.Go(func() error { return afcProcessor.Run(ctx) })
	eg.Go(func() error { return managerScheduler.Run(ctx) })

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}

func initManagerScheduler(
	cfg config.ManagerSchedulerService,
	managerPool managerpool.Pool,
	msgRepo *messagesrepo.Repo,
	outbox *outbox.Service,
	problemsRepo *problemsrepo.Repo,
	txtor *storegen.Database,
	log *zap.Logger,
) (*managerscheduler.Service, error) {
	return managerscheduler.New(
		managerscheduler.NewOptions(cfg.IdleDuration, managerPool, msgRepo, outbox, problemsRepo, txtor, log),
	)
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
	messagesRepo *messagesrepo.Repo,
	chatsRepo *chatsrepo.Repo,
	problemsRepo *problemsrepo.Repo,
	managerLoad *managerload.Service,
	eventStream eventstream.EventStream,
) (*outbox.Service, error) {
	repo, err := jobsrepo.New(jobsrepo.NewOptions(storeDB))
	if err != nil {
		return nil, fmt.Errorf("create jobs repo: %v", err)
	}
	srvLogger := zap.L().Named(logNameOutbox)
	outboxSvc, err := outbox.New(outbox.NewOptions(cfg.Workers, cfg.IdleTime, cfg.ReserveFor, repo, storeDB, srvLogger))
	if err != nil {
		return nil, fmt.Errorf("create outbox service: %v", err)
	}

	jobs, err := createOutboxJobs(msgProducer, messagesRepo, chatsRepo, problemsRepo, managerLoad, eventStream, srvLogger)
	if err != nil {
		return nil, fmt.Errorf("create outbox jobs: %v", err)
	}
	for _, job := range jobs {
		err = outboxSvc.RegisterJob(job)
		if err != nil {
			return nil, fmt.Errorf("register %q job: %v", job.Name(), err)
		}
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

func initWebsocketHandler(
	lg *zap.Logger,
	allowOrigins []string,
	secWsProtocol string,
	eventStream eventstream.EventStream,
	eventsAdepter websocketstream.EventAdapter,
) (
	*websocketstream.HTTPHandler,
	error,
) {
	shutdownCh := make(chan struct{})
	wsHandler, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		lg,
		eventStream,
		eventsAdepter,
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

func initAFCVerdictsProcessor(
	cfg config.AFCVerdictsProcessorService,
	storeDB *storegen.Database,
	messagesRepo *messagesrepo.Repo,
	outboxSvc *outbox.Service,
) (*afcverdictsprocessor.Service, error) {
	return afcverdictsprocessor.New(afcverdictsprocessor.NewOptions(
		cfg.Brokers,
		cfg.Consumers,
		cfg.ConsumerGroup,
		cfg.VerdictsTopic,
		afcverdictsprocessor.NewKafkaReader,
		afcverdictsprocessor.NewKafkaDLQWriter(cfg.Brokers, cfg.VerdictsDLQTopic),
		storeDB,
		messagesRepo,
		outboxSvc,
		zap.L().Named(logNameAFCProcessor),
		afcverdictsprocessor.WithVerdictsSignKey(cfg.VerdictsSignPublicKey),
	))
}

func createOutboxJobs(
	msgProducer *msgproducer.Service,
	messagesRepo *messagesrepo.Repo,
	chatsRepo *chatsrepo.Repo,
	problemsRepo *problemsrepo.Repo,
	managerLoad *managerload.Service,
	eventStream eventstream.EventStream,
	log *zap.Logger,
) ([]outbox.Job, error) {
	jobs := make([]outbox.Job, 0)
	sendClientMsgJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(
		msgProducer,
		messagesRepo,
		eventStream,
		log,
	))
	if err != nil {
		return nil, fmt.Errorf("create send-client-message job: %v", err)
	}
	jobs = append(jobs, sendClientMsgJob)

	clientMessageSentJob, err := clientmessagesentjob.New(
		clientmessagesentjob.NewOptions(messagesRepo, problemsRepo, eventStream, log),
	)
	if err != nil {
		return nil, fmt.Errorf("create %q job: %v", clientmessagesentjob.Name, err)
	}
	jobs = append(jobs, clientMessageSentJob)

	clientMessageBlockedJob, err := clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(messagesRepo, eventStream, log))
	if err != nil {
		return nil, fmt.Errorf("create %q job: %v", clientmessagesentjob.Name, err)
	}
	jobs = append(jobs, clientMessageBlockedJob)

	managerAssignedToProblemJob, err := managerassignedtoproblemjob.New(
		managerassignedtoproblemjob.NewOptions(messagesRepo, problemsRepo, managerLoad, eventStream, log),
	)
	if err != nil {
		return nil, fmt.Errorf("create %q job: %v", managerassignedtoproblemjob.Name, err)
	}
	jobs = append(jobs, managerAssignedToProblemJob)

	sendManagerMsgJob, err := sendmanagermessage.New(
		sendmanagermessage.NewOptions(msgProducer, messagesRepo, chatsRepo, eventStream, log),
	)
	if err != nil {
		return nil, fmt.Errorf("create %q job: %v", sendmanagermessage.Name, err)
	}
	jobs = append(jobs, sendManagerMsgJob)

	closeChatJob, err := closechatjob.New(closechatjob.NewOptions(problemsRepo, messagesRepo, eventStream, log))
	if err != nil {
		return nil, fmt.Errorf("create %q job: %v", closechatjob.Name, err)
	}
	jobs = append(jobs, closeChatJob)

	return jobs, nil
}
