package state

import "sync"

type Machine struct {
	userSteps map[int64]State
	mu        sync.RWMutex
}

func New() *Machine {
	return &Machine{
		userSteps: make(map[int64]State),
	}
}

func (m *Machine) GetStep(userID int64) State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.userSteps[userID]
}

func (m *Machine) SetStep(userID int64, state State) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userSteps[userID] = state
}
