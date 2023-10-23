package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/ekhvalov/bank-chat-service/internal/clients/keycloak"
	"github.com/ekhvalov/bank-chat-service/internal/config"
	chatsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/chats"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	serverclient "github.com/ekhvalov/bank-chat-service/internal/server-client"
	"github.com/ekhvalov/bank-chat-service/internal/server-client/errhandler"
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
	"github.com/ekhvalov/bank-chat-service/internal/store"
	gethistory "github.com/ekhvalov/bank-chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/ekhvalov/bank-chat-service/internal/usecases/client/send-message"
)

const nameServerClient = "server-client"

func initServerClient(
	cfg config.ClientServerConfig,
	v1Swagger *openapi3.T,
	client *keycloakclient.Client,
	msgRepo *messagesrepo.Repo,
	chatsRepo *chatsrepo.Repo,
	problemsRepo *problemsrepo.Repo,
	transactor *store.Database,
	isProduction bool,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("create gethistory usecase: %v", err)
	}

	sendMessageUseCase, err := sendmessage.New(sendmessage.NewOptions(
		chatsRepo,
		msgRepo,
		problemsRepo,
		transactor,
	))
	if err != nil {
		return nil, fmt.Errorf("create sendMessage usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg, getHistoryUseCase, sendMessageUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(lg, isProduction, errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	srv, err := serverclient.New(serverclient.NewOptions(
		cfg.Addr,
		cfg.AllowOrigins,
		cfg.RequiredAccess.Resource,
		cfg.RequiredAccess.Role,
		lg,
		v1Swagger,
		v1Handlers,
		client,
		errHandler.Handle,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
