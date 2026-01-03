package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JSONManager struct {
	configPath string
}

func NewJSONManager(configPath string) *JSONManager {
	return &JSONManager{
		configPath: configPath,
	}
}

func (m *JSONManager) Load() (*Config, error) {
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		return &Config{Connections: []Connection{}}, nil
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Connections == nil {
		config.Connections = []Connection{}
	}

	return &config, nil
}

func (m *JSONManager) Save(config *Config) error {
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

