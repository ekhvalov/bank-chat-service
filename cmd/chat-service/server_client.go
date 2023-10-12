package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/ekhvalov/bank-chat-service/internal/clients/keycloak"
	"github.com/ekhvalov/bank-chat-service/internal/config"
	serverclient "github.com/ekhvalov/bank-chat-service/internal/server-client"
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
)

const nameServerClient = "server-client"

func initServerClient(
	cfg config.ClientServerConfig,
	v1Swagger *openapi3.T,
	client *keycloakclient.Client,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
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
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
