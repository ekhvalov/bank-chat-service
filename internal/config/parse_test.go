package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ekhvalov/bank-chat-service/internal/config"
)

var configExamplePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	configExamplePath = filepath.Join(filepath.Dir(currentFile), "..", "..", "configs", "config.example.toml")
}

func TestParseAndValidate(t *testing.T) {
	cfg, err := config.ParseAndValidate(configExamplePath)
	require.NoError(t, err)
	assert.NotEmpty(t, cfg.Log.Level)
}

func TestParseAndValidateWithEnv(t *testing.T) {
	assert.NoError(t, os.Setenv(fmt.Sprintf("%s_LOG_LEVEL", config.EnvPrefix), "debug"))
	cfg, err := config.ParseAndValidate(configExamplePath)
	require.NoError(t, err)
	assert.NotEmpty(t, cfg.Log.Level)
	assert.Equal(t, "debug", cfg.Log.Level)
}
