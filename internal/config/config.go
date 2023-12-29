package config

import "time"

// Documentation https://pkg.go.dev/github.com/go-playground/validator/v10

type Config struct {
	Global   GlobalConfig   `toml:"global"`
	Log      LogConfig      `toml:"log"`
	Clients  ClientsConfig  `toml:"clients"`
	Servers  ServersConfig  `toml:"servers"`
	Services ServicesConfig `toml:"services"`
	Sentry   SentryConfig   `toml:"sentry"`
}

type GlobalConfig struct {
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

func (gc GlobalConfig) IsProduction() bool {
	return gc.Env == "prod"
}

type LogConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type ClientsConfig struct {
	Keycloak KeycloakClientConfig `toml:"keycloak" validate:"required"`
	Postgres PostgresClientConfig `toml:"postgres" validate:"required"`
}

type KeycloakClientConfig struct {
	BasePath     string `toml:"base_path" validate:"required,url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode"`
}

type PostgresClientConfig struct {
	Address   string `toml:"address" validate:"required,hostname_port"`
	Username  string `toml:"username" validate:"required"`
	Password  string `toml:"password" validate:"required"`
	Database  string `toml:"database" validate:"required"`
	DebugMode bool   `toml:"debug_mode"`
}

type ServersConfig struct {
	Client  ClientServerConfig  `toml:"client"`
	Manager ManagerServerConfig `toml:"manager"`
	Debug   DebugServerConfig   `toml:"debug"`
}

type ClientServerConfig struct {
	Addr           string               `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string             `toml:"allow_origins" validate:"dive,required,http_url"`
	RequiredAccess RequiredAccessConfig `toml:"required_access"`
	SecWsProtocol  string               `toml:"sec_ws_protocol" validate:"required"`
}

type ManagerServerConfig struct {
	Addr           string               `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string             `toml:"allow_origins" validate:"dive,required,http_url"`
	RequiredAccess RequiredAccessConfig `toml:"required_access"`
	SecWsProtocol  string               `toml:"sec_ws_protocol" validate:"required"`
}

type RequiredAccessConfig struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type ServicesConfig struct {
	MsgProducer                 MsgProducerServiceConfig    `toml:"msg_producer"`
	ManagerLoad                 ManagerLoadService          `toml:"manager_load"`
	OutboxService               OutboxService               `toml:"outbox"`
	AFCVerdictsProcessorService AFCVerdictsProcessorService `toml:"afc_verdicts_processor"`
	ManagerSchedulerService     ManagerSchedulerService     `toml:"manager_scheduler"`
}

type MsgProducerServiceConfig struct {
	Brokers    []string `toml:"brokers" validate:"dive,required,hostname_port"`
	Topic      string   `toml:"topic" validate:"required"`
	BatchSize  int      `toml:"batch_size" validate:"min=1"`
	EncryptKey string   `toml:"encrypt_key" validate:"omitempty,hexadecimal"`
}

type ManagerLoadService struct {
	MaxProblemsAtSameTime int `toml:"max_problems_at_same_time" validate:"min=1"`
}

type ManagerSchedulerService struct {
	IdleDuration time.Duration `toml:"idle_duration" validate:"min=100ms,max=1m"`
}

type OutboxService struct {
	Workers    int           `toml:"workers" validate:"min=1"`
	IdleTime   time.Duration `toml:"idle_time" validate:"required"`
	ReserveFor time.Duration `toml:"reserve_for" validate:"required"`
}

type AFCVerdictsProcessorService struct {
	Brokers               []string `toml:"brokers" validate:"dive,required,hostname_port"`
	Consumers             int      `toml:"consumers" validate:"min=1,max=16"`
	ConsumerGroup         string   `toml:"consumer_group" validate:"required"`
	VerdictsTopic         string   `toml:"verdicts_topic" validate:"required"`
	VerdictsSignPublicKey string   `toml:"verdicts_signing_public_key" validate:"omitempty,min=1"`
	VerdictsDLQTopic      string   `toml:"verdicts_dlq_topic" validate:"required"`
}

type SentryConfig struct {
	DSN string `toml:"dsn" validate:"omitempty,http_url"`
}
