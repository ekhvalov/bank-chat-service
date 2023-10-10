package config

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/ekhvalov/bank-chat-service/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return Config{}, fmt.Errorf("decode config: %v", err)
	}

	if err := validator.Validator.Struct(config); err != nil {
		return Config{}, fmt.Errorf("struct validate: %v", err)
	}
	return config, nil
}
