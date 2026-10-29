package storage

import "sync"

type Memory struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewMemory() *Memory {
	return &Memory{
		links: make(map[string]string),
	}
}

func (m *Memory) Save(code, url string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.links[code] = url
}

func (m *Memory) Get(code string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, ok := m.links[code]
	return url, ok
}
