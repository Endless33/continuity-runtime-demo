package session

import "fmt"

type State string

const (
	StateInit        State = "INIT"
	StateEstablished State = "ESTABLISHED"
)

type Session struct {
	ID       string
	Epoch    int
	State    State
	LastAddr string
}

func NewSession(id string) *Session {
	return &Session{
		ID:    id,
		Epoch: 1,
		State: StateInit,
	}
}

func (s *Session) UpdatePath(addr string) {
	if s.LastAddr != "" && s.LastAddr != addr {
		fmt.Println("path change:", s.LastAddr, "->", addr)
	}
	s.LastAddr = addr
}

func (s *Session) Establish() {
	if s.State != StateEstablished {
		s.State = StateEstablished
		fmt.Println("session established:", s.ID, "epoch:", s.Epoch)
	}
}