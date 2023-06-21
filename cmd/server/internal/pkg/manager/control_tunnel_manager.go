package manager

import (
	"sync"

	"github.com/3vilive/expit/cmd/server/internal/pkg/tunnel"
)

var gControlTunnelManager *ControlTunnelManager

func init() {
	gControlTunnelManager = &ControlTunnelManager{
		controlTunnelMap: make(map[string]*tunnel.ControlTunnel),
	}
}

func GetControlTunnelManager() *ControlTunnelManager {
	return gControlTunnelManager
}

type ControlTunnelManager struct {
	mutex            sync.RWMutex
	controlTunnelMap map[string]*tunnel.ControlTunnel
}

func (m *ControlTunnelManager) GetControlTunnelById(id string) *tunnel.ControlTunnel {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.controlTunnelMap[id]
}

func (m *ControlTunnelManager) RemoveControlTunnelById(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.controlTunnelMap, id)
}

func (m *ControlTunnelManager) AddControlTunnel(id string, t *tunnel.ControlTunnel) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.controlTunnelMap[id] = t
}
