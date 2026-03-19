package runtime

import (
	"fmt"
)

type Session struct {
	ID            string
	Epoch         int
	ActivePath    string
	ActiveTransport string
	Alive         bool
}

func NewSession(id string, transport Transport) *Session {
	return &Session{
		ID:              id,
		Epoch:           1,
		ActivePath:      transport.Name,
		ActiveTransport: transport.Name,
		Alive:           true,
	}
}

func (s *Session) TransferAuthority(newTransport Transport) {
	s.Epoch++

	fmt.Printf("[SESSION] authority transfer → %s (epoch=%d)\n",
		newTransport.Name,
		s.Epoch,
	)

	s.ActiveTransport = newTransport.Name
	s.ActivePath = newTransport.Name
}

func (s *Session) ValidateTransport(name string, epoch int) bool {
	if !s.Alive {
		return false
	}

	if epoch != s.Epoch {
		fmt.Printf("[SESSION] reject transport=%s (stale epoch %d)\n", name, epoch)
		return false
	}

	if name != s.ActiveTransport {
		fmt.Printf("[SESSION] reject transport=%s (not active)\n", name)
		return false
	}

	return true
}