package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

func ParseAndValidate(filename string) (Config, error) {
	config, parseErr := parse(filename)
	if parseErr != nil {
		return Config{}, parseErr
	}
	if validateErr := validate(config); validateErr != nil {
		return Config{}, validateErr
	}
	return config, nil
}

func parse(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("read file error: %v", err)
	}
	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("config unmarshal error: %v", err)
	}
	return config, nil
}

func validate(config Config) error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(config)
}
