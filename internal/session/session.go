package session

import (
	"fmt"
	"sync"
)

type Session struct {
	ID        string
	Epoch     int
	LastAddr  string
	mu        sync.Mutex
}

func NewSession(id string) *Session {
	return &Session{
		ID:    id,
		Epoch: 1,
	}
}

func (s *Session) UpdatePath(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.LastAddr != addr {
		fmt.Println("🔁 path change:", s.LastAddr, "→", addr)
		s.LastAddr = addr
	}
}

func (s *Session) NextEpoch() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Epoch++
	fmt.Println("⬆ epoch:", s.Epoch)
}