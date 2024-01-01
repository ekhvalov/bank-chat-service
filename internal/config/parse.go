package config

import (
	"fmt"
	"path"

	"github.com/kkyr/fig"
)

const EnvPrefix = "BANKCHAT"

func ParseAndValidate(filepath string) (Config, error) {
	var config Config

	err := fig.Load(&config, fig.UseEnv(EnvPrefix), fig.File(path.Base(filepath)), fig.Dirs(path.Dir(filepath)))
	if err != nil {
		return Config{}, fmt.Errorf("load config: %v", err)
	}

	return config, nil
}
