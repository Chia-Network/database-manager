package manager

import (
	"github.com/chia-network/database-manager/internal/config"
	"github.com/chia-network/database-manager/internal/iface"
)

// Manager is the database manager tool that connects to the selected database engine and applies the config
type Manager struct {
	manager iface.Manager
	cfg     *config.Config
}

// NewManager returns a new instance of the Manager with the specified DB engine
func NewManager(specificManager iface.Manager, cfg *config.Config) (*Manager, error) {
	return &Manager{
		manager: specificManager,
		cfg:     cfg,
	}, nil
}

// Apply Applies the configuration to the database
func (m *Manager) Apply() error {
	var err error
	for _, database := range m.cfg.Databases {
		err = m.manager.CreateDatabase(database.Name)
		if err != nil {
			return err
		}
	}
	// @TODO
	return nil
}
