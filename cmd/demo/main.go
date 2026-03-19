package main

import (
	"fmt"
	"time"
)

type State string
type Event string

const (
	StateAttached   State = "ATTACHED"
	StateRecovering State = "RECOVERING"

	EventWiFiFailed Event = "WiFi failed"
)

type Packet struct {
	ID int
}

type Transport struct {
	Name    string
	Latency time.Duration
	Score   float64
}

type Runtime struct {
	state      State
	epoch      int
	current    Transport
	candidates []Transport
	packetID   int
}

func NewRuntime() *Runtime {
	return &Runtime{
		state: StateAttached,
		epoch: 1,
		current: Transport{
			Name:    "wifi",
			Latency: 120 * time.Millisecond,
			Score:   20,
		},
		candidates: []Transport{
			{Name: "5g", Latency: 40 * time.Millisecond, Score: 100},
			{Name: "lte", Latency: 80 * time.Millisecond, Score: 60},
		},
		packetID: 100,
	}
}

func (r *Runtime) sendPacket() {
	r.packetID++
	fmt.Printf("[SEND] packet #%d via %s\n", r.packetID, r.current.Name)
	time.Sleep(r.current.Latency)
}

func (r *Runtime) HandleEvent(e Event) {
	switch e {
	case EventWiFiFailed:
		r.onWiFiFailed()
	}
}

func (r *Runtime) onWiFiFailed() {
	fmt.Println("\n[EVENT] WiFi failed")

	best := r.selectBestTransport()

	fmt.Printf("[DECISION] best=%s (score=%.1f, latency=%v)\n",
		best.Name, best.Score, best.Latency)

	if best.Score > r.current.Score {
		fmt.Println("[DECISION] migrate=true")

		r.transition(StateRecovering)

		time.Sleep(300 * time.Millisecond)

		r.epoch++
		fmt.Printf("[AUTHORITY] epoch %d granted to %s\n", r.epoch, best.Name)

		fmt.Printf("[CHECK] stale %s rejected\n", r.current.Name)

		r.current = best

		r.transition(StateAttached)

		fmt.Println("[RESULT] session continues (no reconnect)")
	} else {
		fmt.Println("[RESULT] no better transport")
	}
}

func (r *Runtime) selectBestTransport() Transport {
	best := r.candidates[0]

	for _, t := range r.candidates {
		if t.Score > best.Score {
			best = t
		}
	}

	return best
}

func (r *Runtime) transition(newState State) {
	fmt.Printf("[STATE] %s -> %s\n", r.state, newState)
	r.state = newState
}

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (PACKET FLOW + LATENCY)")
	fmt.Println()

	r := NewRuntime()

	// send packets before failure
	for i := 0; i < 3; i++ {
		r.sendPacket()
	}

	// simulate failure
	r.HandleEvent(EventWiFiFailed)

	// continue sending after migration
	for i := 0; i < 3; i++ {
		r.sendPacket()
	}
}