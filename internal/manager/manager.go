package manager

import (
	"github.com/chia-network/database-manager/internal/iface"
)

// Manager is the database manager tool that connects to the selected database engine and applies the config
type Manager struct {
	manager iface.Manager
}

// NewManager returns a new instance of the Manager with the specified DB engine
func NewManager(specificManager iface.Manager) (*Manager, error) {
	return &Manager{manager: specificManager}, nil
}

// Apply Applies the configuration to the database
func (m *Manager) Apply() error {
	// @TODO
	return nil
}
