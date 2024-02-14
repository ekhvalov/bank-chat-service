package config

import "time"

// Documentation https://pkg.go.dev/github.com/go-playground/validator/v10

type Config struct {
	Global   GlobalConfig   `fig:"global"`
	Log      LogConfig      `fig:"log"`
	Clients  ClientsConfig  `fig:"clients"`
	Servers  ServersConfig  `fig:"servers"`
	Services ServicesConfig `fig:"services"`
	Sentry   SentryConfig   `fig:"sentry"`
}

type GlobalConfig struct {
	Env string `fig:"env" validate:"required,oneof=dev stage prod"`
}

func (gc GlobalConfig) IsProduction() bool {
	return gc.Env == "prod"
}

type LogConfig struct {
	Level string `fig:"level" validate:"required,oneof=debug info warn error"`
}

type ClientsConfig struct {
	Keycloak KeycloakClientConfig `fig:"keycloak" validate:"required"`
	Postgres PostgresClientConfig `fig:"postgres" validate:"required"`
}

type KeycloakClientConfig struct {
	BasePath     string `fig:"base_path" validate:"required,url"`
	Realm        string `fig:"realm" validate:"required"`
	ClientID     string `fig:"client_id" validate:"required"`
	ClientSecret string `fig:"client_secret" validate:"required"`
	DebugMode    bool   `fig:"debug_mode"`
}

type PostgresClientConfig struct {
	Address   string `fig:"address" validate:"required,hostname_port"`
	Username  string `fig:"username" validate:"required"`
	Password  string `fig:"password" validate:"required"`
	Database  string `fig:"database" validate:"required"`
	DebugMode bool   `fig:"debug_mode"`
}

type ServersConfig struct {
	Client  ClientServerConfig  `fig:"client"`
	Manager ManagerServerConfig `fig:"manager"`
	Debug   DebugServerConfig   `fig:"debug"`
}

type ClientServerConfig struct {
	Addr           string               `fig:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string             `fig:"allow_origins" validate:"dive,required,http_url"`
	RequiredAccess RequiredAccessConfig `fig:"required_access"`
	SecWsProtocol  string               `fig:"sec_ws_protocol" validate:"required"`
}

type ManagerServerConfig struct {
	Addr           string               `fig:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string             `fig:"allow_origins" validate:"dive,required,http_url"`
	RequiredAccess RequiredAccessConfig `fig:"required_access"`
	SecWsProtocol  string               `fig:"sec_ws_protocol" validate:"required"`
}

type RequiredAccessConfig struct {
	Resource string `fig:"resource" validate:"required"`
	Role     string `fig:"role" validate:"required"`
}

type DebugServerConfig struct {
	Addr string `fig:"addr" validate:"required,hostname_port"`
}

type ServicesConfig struct {
	MsgProducer                 MsgProducerServiceConfig    `fig:"msg_producer"`
	ManagerLoad                 ManagerLoadService          `fig:"manager_load"`
	OutboxService               OutboxService               `fig:"outbox"`
	AFCVerdictsProcessorService AFCVerdictsProcessorService `fig:"afc_verdicts_processor"`
	ManagerSchedulerService     ManagerSchedulerService     `fig:"manager_scheduler"`
}

type MsgProducerServiceConfig struct {
	Brokers    []string `fig:"brokers" validate:"dive,required,hostname_port"`
	Topic      string   `fig:"topic" validate:"required"`
	BatchSize  int      `fig:"batch_size" validate:"min=1"`
	EncryptKey string   `fig:"encrypt_key" validate:"omitempty,hexadecimal"`
}

type ManagerLoadService struct {
	MaxProblemsAtSameTime int `fig:"max_problems_at_same_time" validate:"min=1"`
}

type ManagerSchedulerService struct {
	IdleDuration time.Duration `fig:"idle_duration" validate:"min=100ms,max=1m"`
}

type OutboxService struct {
	Workers    int           `fig:"workers" validate:"min=1"`
	IdleTime   time.Duration `fig:"idle_time" validate:"required"`
	ReserveFor time.Duration `fig:"reserve_for" validate:"required"`
}

type AFCVerdictsProcessorService struct {
	Brokers               []string `fig:"brokers" validate:"dive,required,hostname_port"`
	Consumers             int      `fig:"consumers" validate:"min=1,max=16"`
	ConsumerGroup         string   `fig:"consumer_group" validate:"required"`
	VerdictsTopic         string   `fig:"verdicts_topic" validate:"required"`
	VerdictsSignPublicKey string   `fig:"verdicts_signing_public_key" validate:"omitempty,min=1"`
	VerdictsDLQTopic      string   `fig:"verdicts_dlq_topic" validate:"required"`
}

type SentryConfig struct {
	DSN string `fig:"dsn" validate:"omitempty,http_url"`
}
