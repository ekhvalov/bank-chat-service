package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/config"
	chatsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/chats"
	messagesrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/messages"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	"github.com/ekhvalov/bank-chat-service/internal/server"
	servermanager "github.com/ekhvalov/bank-chat-service/internal/server-manager"
	errhandlermanager "github.com/ekhvalov/bank-chat-service/internal/server-manager/errhandler"
	managerevents "github.com/ekhvalov/bank-chat-service/internal/server-manager/events"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
	"github.com/ekhvalov/bank-chat-service/internal/server/errhandler"
	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	managerload "github.com/ekhvalov/bank-chat-service/internal/services/manager-load"
	managerpool "github.com/ekhvalov/bank-chat-service/internal/services/manager-pool"
	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chats"
)

const (
	logNameServerManager = "server-manager"
)

func initServerManager(
	cfg config.Config,
	eventStream eventstream.EventStream,
	managerPool managerpool.Pool,
	chatsRepo *chatsrepo.Repo,
	messagesRepo *messagesrepo.Repo,
	problemsRepo *problemsrepo.Repo,
) (*server.Server, error) {
	lg := zap.L().Named(logNameServerManager)

	v1Swagger, err := managerv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get swagger: %v", err)
	}

	v1Handlers, err := initV1ManagerHandlers(lg, cfg.Services.ManagerLoad, managerPool, chatsRepo, messagesRepo, problemsRepo)
	if err != nil {
		return nil, fmt.Errorf("create v1Handlers: %v", err)
	}

	handlersRegistrar, err := servermanager.New(servermanager.NewOptions(v1Swagger, v1Handlers))
	if err != nil {
		return nil, fmt.Errorf("create handlers registrar: %v", err)
	}

	keycloakClient, err := initKeycloakClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create keycloak client: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(lg, cfg.Global.IsProduction(), errhandlermanager.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	wsHandler, err := initWebsocketHandler(
		lg,
		cfg.Servers.Manager.AllowOrigins,
		cfg.Servers.Manager.SecWsProtocol,
		eventStream, managerevents.Adapter{},
	)
	if err != nil {
		return nil, fmt.Errorf("create websocket: %v", err)
	}

	s, err := server.New(server.NewOptions(
		cfg.Servers.Manager.Addr,
		cfg.Servers.Manager.AllowOrigins,
		cfg.Servers.Manager.RequiredAccess.Resource,
		cfg.Servers.Manager.RequiredAccess.Role,
		cfg.Servers.Manager.SecWsProtocol,
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

func initV1ManagerHandlers(
	lg *zap.Logger,
	cfg config.ManagerLoadService,
	managerPool managerpool.Pool,
	chatsRepo *chatsrepo.Repo,
	messagesRepo *messagesrepo.Repo,
	problemsRepo *problemsrepo.Repo,
) (managerv1.Handlers, error) {
	loadService, err := managerload.New(managerload.NewOptions(cfg.MaxProblemsAtSameTime, problemsRepo))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create manager load service: %v", err)
	}

	canReceiveProblemsUsecase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(loadService, managerPool))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create canReceiveProblems usecase: %v", err)
	}

	freeHandsUsecase, err := freehands.New(freehands.NewOptions(loadService, managerPool))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create freeHands usecase: %v", err)
	}

	getChatsUsecase, err := getchats.New(getchats.NewOptions(chatsRepo))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create getChats usecase: %v", err)
	}

	getChatHistoryUsecase, err := getchathistory.New(getchathistory.NewOptions(messagesRepo, problemsRepo))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create getChatHistoryUsecase: %v", err)
	}

	return managerv1.NewHandlers(
		managerv1.NewOptions(lg, canReceiveProblemsUsecase, freeHandsUsecase, getChatsUsecase, getChatHistoryUsecase),
	)
}
