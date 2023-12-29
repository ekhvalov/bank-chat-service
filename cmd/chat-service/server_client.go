package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/buildinfo"
	keycloakclient "github.com/ekhvalov/bank-chat-service/internal/clients/keycloak"
	"github.com/ekhvalov/bank-chat-service/internal/config"
	chatsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/chats"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	"github.com/ekhvalov/bank-chat-service/internal/server"
	serverclient "github.com/ekhvalov/bank-chat-service/internal/server-client"
	errhandlerclient "github.com/ekhvalov/bank-chat-service/internal/server-client/errhandler"
	clientevents "github.com/ekhvalov/bank-chat-service/internal/server-client/events"
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
	"github.com/ekhvalov/bank-chat-service/internal/server/errhandler"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	"github.com/ekhvalov/bank-chat-service/internal/services/outbox"
	storegen "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	gethistory "github.com/ekhvalov/bank-chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/ekhvalov/bank-chat-service/internal/usecases/client/send-message"
)

const (
	logNameServerClient = "server-client"
)

func initServerClient(
	cfg config.Config,
	storeDB *storegen.Database,
	outboxSvc *outbox.Service,
	eventStream eventstream.EventStream,
) (*server.Server, error) {
	lg := zap.L().Named(logNameServerClient)

	v1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get swagger: %v", err)
	}
	v1Handlers, err := initV1ClientHandlers(lg, storeDB, outboxSvc)
	if err != nil {
		return nil, fmt.Errorf("create v1Handlers: %v", err)
	}
	handlersRegistrar, err := serverclient.New(serverclient.NewOptions(v1Swagger, v1Handlers))
	if err != nil {
		return nil, fmt.Errorf("create handlers registrar: %v", err)
	}

	keycloakClient, err := initKeycloakClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create keycloak client: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(lg, cfg.Global.IsProduction(), errhandlerclient.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	wsHandler, err := initWebsocketHandler(
		lg,
		cfg.Servers.Client.AllowOrigins,
		cfg.Servers.Client.SecWsProtocol,
		eventStream,
		clientevents.Adapter{},
	)
	if err != nil {
		return nil, fmt.Errorf("create websocket: %v", err)
	}

	s, err := server.New(server.NewOptions(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
		cfg.Servers.Client.SecWsProtocol,
		lg,
		keycloakClient,
		errHandler.Handle,
		wsHandler,
	), handlersRegistrar)
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return s, nil
}

func initKeycloakClient(cfg config.Config) (*keycloakclient.Client, error) {
	return keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode),
		keycloakclient.WithUserAgent(fmt.Sprintf("chat-service/%s", buildinfo.BuildInfo.Main.Version)),
	))
}

func initV1ClientHandlers(lg *zap.Logger, storeDB *storegen.Database, outboxSvc *outbox.Service) (clientv1.Handlers, error) {
	messagesRepo, err := messagesrepo.New(messagesrepo.NewOptions(storeDB))
	if err != nil {
		return clientv1.Handlers{}, fmt.Errorf("create messages repo: %v", err)
	}

	chatsRepo, err := chatsrepo.New(chatsrepo.NewOptions(storeDB))
	if err != nil {
		return clientv1.Handlers{}, fmt.Errorf("create chats repo: %v", err)
	}

	problemsRepo, err := problemsrepo.New(problemsrepo.NewOptions(storeDB))
	if err != nil {
		return clientv1.Handlers{}, fmt.Errorf("create problems repo: %v", err)
	}

	sendMessageOptions := sendmessage.NewOptions(chatsRepo, messagesRepo, outboxSvc, problemsRepo, storeDB)
	sendMessageUseCase, err := sendmessage.New(sendMessageOptions)
	if err != nil {
		return clientv1.Handlers{}, fmt.Errorf("create sendmessage usecase: %v", err)
	}

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(messagesRepo))
	if err != nil {
		return clientv1.Handlers{}, fmt.Errorf("create gethistory usecase: %v", err)
	}

	return clientv1.NewHandlers(clientv1.NewOptions(lg, getHistoryUseCase, sendMessageUseCase))
}
