package runtime

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

type Runtime struct {
	State      State
	Epoch      int
	Current    Transport
	Candidates []Transport
	PacketID   int
}

func NewRuntime() *Runtime {
	return &Runtime{
		State: StateAttached,
		Epoch: 1,
		Current: Transport{
			Name:    "wifi",
			Latency: 120 * time.Millisecond,
			Score:   20,
		},
		Candidates: []Transport{
			{Name: "5g", Latency: 40 * time.Millisecond, Score: 100},
			{Name: "lte", Latency: 80 * time.Millisecond, Score: 60},
		},
		PacketID: 100,
	}
}

func (r *Runtime) SendPacket() {
	r.PacketID++
	fmt.Printf("[SEND] packet #%d via %s\n", r.PacketID, r.Current.Name)
	time.Sleep(r.Current.Latency)
}

func (r *Runtime) HandleEvent(e Event) {
	switch e {
	case EventWiFiFailed:
		r.onWiFiFailed()
	}
}

func (r *Runtime) onWiFiFailed() {
	fmt.Println("\n[EVENT] WiFi failed")

	best := SelectBestTransport(r.Current, r.Candidates)

	fmt.Printf("[DECISION] best=%s (score=%.1f)\n", best.Name, best.Score)

	if best.Score > r.Current.Score {
		fmt.Println("[DECISION] migrate=true")

		r.transition(StateRecovering)

		time.Sleep(300 * time.Millisecond)

		r.Epoch++
		fmt.Printf("[AUTHORITY] epoch %d granted to %s\n", r.Epoch, best.Name)

		fmt.Printf("[CHECK] stale %s rejected\n", r.Current.Name)

		r.Current = best

		r.transition(StateAttached)

		fmt.Println("[RESULT] session continues (no reconnect)")
	} else {
		fmt.Println("[RESULT] no better transport")
	}
}

func (r *Runtime) transition(newState State) {
	fmt.Printf("[STATE] %s -> %s\n", r.State, newState)
	r.State = newState
}