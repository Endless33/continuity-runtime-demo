package runtime

import (
	"fmt"
	"math/rand"
	"time"
)

type NetworkSimulator struct {
	LossRate float64
	Jitter   time.Duration
}

func NewNetworkSimulator() *NetworkSimulator {
	return &NetworkSimulator{
		LossRate: 0.2,                  // 20% packet loss
		Jitter:   50 * time.Millisecond,
	}
}

func (n *NetworkSimulator) Transmit(packetID int, t Transport) bool {
	// simulate packet loss
	if rand.Float64() < n.LossRate {
		fmt.Printf("[LOSS] packet #%d dropped on %s\n", packetID, t.Name)
		return false
	}

	// simulate latency + jitter
	delay := t.Latency + time.Duration(rand.Int63n(int64(n.Jitter)))
	time.Sleep(delay)

	fmt.Printf("[RECV] packet #%d via %s (delay=%v)\n", packetID, t.Name, delay)
	return true
}