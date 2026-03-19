package runtime

import "fmt"

type ACKTracker struct {
	acked map[int]bool
}

func NewACKTracker() *ACKTracker {
	return &ACKTracker{
		acked: make(map[int]bool),
	}
}

func (a *ACKTracker) Ack(id int) {
	a.acked[id] = true
	fmt.Printf("[ACK] packet #%d acknowledged\n", id)
}

func (a *ACKTracker) IsAcked(id int) bool {
	return a.acked[id]
}