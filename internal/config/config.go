package config

// Documentation https://pkg.go.dev/github.com/go-playground/validator/v10

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Clients ClientsConfig `toml:"clients"`
	Servers ServersConfig `toml:"servers"`
	Sentry  SentryConfig  `toml:"sentry"`
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
}

type KeycloakClientConfig struct {
	BasePath     string `toml:"base_path" validate:"required,url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode"`
}

type ServersConfig struct {
	Client ClientServerConfig `toml:"client"`
	Debug  DebugServerConfig  `toml:"debug"`
}

type ClientServerConfig struct {
	Addr           string               `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string             `toml:"allow_origins" validate:"dive,required,http_url"`
	RequiredAccess RequiredAccessConfig `toml:"required_access"`
}

type RequiredAccessConfig struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type SentryConfig struct {
	DSN string `toml:"dsn" validate:"omitempty,http_url"`
}
