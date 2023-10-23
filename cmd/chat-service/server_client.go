package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/buildinfo"
	keycloakclient "github.com/ekhvalov/bank-chat-service/internal/clients/keycloak"
	"github.com/ekhvalov/bank-chat-service/internal/config"
	chatsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/chats"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	serverclient "github.com/ekhvalov/bank-chat-service/internal/server-client"
	"github.com/ekhvalov/bank-chat-service/internal/server-client/errhandler"
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
	"github.com/ekhvalov/bank-chat-service/internal/store"
	storegen "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	gethistory "github.com/ekhvalov/bank-chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/ekhvalov/bank-chat-service/internal/usecases/client/send-message"
)

const nameServerClient = "server-client"

func initServerClient(ctx context.Context, cfg config.Config) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	v1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get swagger: %v", err)
	}

	keycloakClient, err := initKeycloakClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create keycloak client: %v", err)
	}

	storeDB, err := initStoreDB(ctx, cfg.Clients.Postgres)
	if err != nil {
		return nil, fmt.Errorf("create store DB: %v", err)
	}

	v1Handlers, err := initV1Handlers(lg, storeDB)
	if err != nil {
		return nil, fmt.Errorf("create v1Handlers: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(lg, cfg.Global.IsProduction(), errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	server, err := serverclient.New(serverclient.NewOptions(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
		lg,
		v1Swagger,
		v1Handlers,
		keycloakClient,
		errHandler.Handle,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return server, nil
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

func initStoreDB(ctx context.Context, cfg config.PostgresClientConfig) (*storegen.Database, error) {
	client, err := initStoreClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create store client: %v", err)
	}
	return storegen.NewDatabase(client), nil
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

	return client, nil
}

func initV1Handlers(lg *zap.Logger, storeDB *storegen.Database) (clientv1.Handlers, error) {
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

	sendMessageOptions := sendmessage.NewOptions(chatsRepo, messagesRepo, problemsRepo, storeDB)
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
