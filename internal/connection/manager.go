package connection

import "dbear/internal/config"

type Manager struct {
	configManager config.Manager
}

func NewManager(configManager config.Manager) *Manager {
	return &Manager{
		configManager: configManager,
	}
}

func (m *Manager) Create(conn config.Connection) error {
	cfg, err := m.configManager.Load()
	if err != nil {
		return err
	}

	for _, existing := range cfg.Connections {
		if existing.Name == conn.Name {
			return nil
		}
	}

	cfg.Connections = append(cfg.Connections, conn)
	return m.configManager.Save(cfg)
}

func (m *Manager) List() ([]config.Connection, error) {
	cfg, err := m.configManager.Load()
	if err != nil {
		return nil, err
	}

	return cfg.Connections, nil
}

func (m *Manager) Get(name string) (*config.Connection, error) {
	cfg, err := m.configManager.Load()
	if err != nil {
		return nil, err
	}

	for _, conn := range cfg.Connections {
		if conn.Name == name {
			return &conn, nil
		}
	}

	return nil, nil
}

func (m *Manager) Delete(name string) error {
	cfg, err := m.configManager.Load()
	if err != nil {
		return err
	}

	filtered := []config.Connection{}
	for _, conn := range cfg.Connections {
		if conn.Name != name {
			filtered = append(filtered, conn)
		}
	}

	cfg.Connections = filtered
	return m.configManager.Save(cfg)
}

