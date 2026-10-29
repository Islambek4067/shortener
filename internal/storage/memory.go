package storage

import (
	"sync"
	"time"
)

type Link struct {
	Original  string
	Clicks    int
	CreatedAt time.Time
}

type Memory struct {
	mu    sync.RWMutex
	links map[string]*Link
}

func NewMemory() *Memory {
	return &Memory{
		links: make(map[string]*Link),
	}
}

func (m *Memory) Save(code, url string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.links[code] = &Link{Original: url, CreatedAt: time.Now()}
}

func (m *Memory) Get(code string) (Link, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.links[code]
	if ok {
		return *l, true
	}
	return Link{}, false
}

func (m *Memory) AddClick(code string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.links[code]
	if ok {
		l.Clicks++
	}
}
