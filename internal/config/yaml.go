package config

type YAMLManager struct {
	configPath string
}

func NewYAMLManager(configPath string) *YAMLManager {
	return &YAMLManager{
		configPath: configPath,
	}
}

func (m *YAMLManager) Load() (*Config, error) {
	return nil, nil
}

func (m *YAMLManager) Save(config *Config) error {
	return nil
}

