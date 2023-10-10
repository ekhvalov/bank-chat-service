package config

// Documentation https://pkg.go.dev/github.com/go-playground/validator/v10

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
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

type ServersConfig struct {
	Client ClientServerConfig `toml:"client"`
	Debug  DebugServerConfig  `toml:"debug"`
}

type ClientServerConfig struct {
	Addr         string   `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins []string `toml:"allow_origins" validate:"dive,required,url"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type SentryConfig struct {
	DSN string `toml:"dsn" validate:"omitempty,http_url"`
}
