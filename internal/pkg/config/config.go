package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HT4w5/bgpsim-go/internal/pkg/model"
)

// Primary config struct
type Config struct {
	Name        string                   `json:"name"`
	Log         *LogConfig               `json:"log"`
	Nodes       []*model.NodeConfig      `json:"nodes"`
	BgpTopology *model.BgpTopologyConfig `json:"bgpTopology"`
}

type LogConfig struct {
	Output string `json:"output"` // Print to console if empty
}

// Create empty config
func New() *Config {
	return &Config{
		Log:   &LogConfig{},
		Nodes: []*model.NodeConfig{},
		BgpTopology: &model.BgpTopologyConfig{
			Nodes: []*model.BgpNodeConfig{},
			Edges: []*model.BgpEdgeConfig{},
		},
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
