package main

import "fmt"

type State string
type Event string

const (
	StateAttached   State = "ATTACHED"
	StateRecovering State = "RECOVERING"

	EventWiFiFailed Event = "WiFi failed"
)

type Transport struct {
	Name  string
	Score float64
}

type Runtime struct {
	state      State
	epoch      int
	current    Transport
	candidates []Transport
}

func NewRuntime() *Runtime {
	return &Runtime{
		state: StateAttached,
		epoch: 1,
		current: Transport{
			Name:  "wifi",
			Score: 20,
		},
		candidates: []Transport{
			{Name: "5g", Score: 100},
			{Name: "lte", Score: 60},
		},
	}
}

func (r *Runtime) HandleEvent(e Event) {
	switch e {
	case EventWiFiFailed:
		r.onWiFiFailed()
	}
}

func (r *Runtime) onWiFiFailed() {
	fmt.Println("[EVENT] WiFi failed")

	best := r.selectBestTransport()

	fmt.Printf("[DECISION] best candidate: %s (score=%.1f)\n", best.Name, best.Score)

	if best.Score > r.current.Score {
		fmt.Println("[DECISION] migrate=true")

		r.transition(StateRecovering)

		r.epoch++
		fmt.Printf("[AUTHORITY] epoch %d granted to %s\n", r.epoch, best.Name)

		fmt.Printf("[CHECK] stale %s rejected\n", r.current.Name)

		r.current = best

		r.transition(StateAttached)

		fmt.Println("[RESULT] session continues (no reconnect)")
	} else {
		fmt.Println("[RESULT] no better transport found")
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
	fmt.Println("CONTINUITY RUNTIME DEMO (SCORING + MIGRATION)")
	fmt.Println()

	r := NewRuntime()

	r.HandleEvent(EventWiFiFailed)
}