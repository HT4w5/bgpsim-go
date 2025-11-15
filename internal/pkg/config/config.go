package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HT4w5/bgpsim-go/internal/pkg/model"
)

// Primary config struct
type Config struct {
	Log     *LogConfig           `json:"log"`
	Network *model.NetworkConfig `json:"network"`
}

type LogConfig struct {
	Output string `json:"output"` // Print to console if empty
}

// Create empty config
func New() *Config {
	return &Config{
		Log:     &LogConfig{},
		Network: &model.NetworkConfig{},
	}
}

// Load config from JSON file
func (c *Config) Load(path string) error {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}

	if err := json.Unmarshal(configBytes, c); err != nil {
		return fmt.Errorf("failed to unmarshal json config: %w", err)
	}

	return nil
}
