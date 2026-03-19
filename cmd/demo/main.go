package main

import "fmt"

type State string
type Event string

const (
	// States
	StateAttached   State = "ATTACHED"
	StateRecovering State = "RECOVERING"

	// Events
	EventWiFiFailed Event = "WiFi failed"
)

// Runtime represents a simple state machine
type Runtime struct {
	state State
	epoch int
}

func NewRuntime() *Runtime {
	return &Runtime{
		state: StateAttached,
		epoch: 1,
	}
}

func (r *Runtime) HandleEvent(e Event) {
	switch e {

	case EventWiFiFailed:
		r.onWiFiFailed()

	default:
		fmt.Println("[WARN] unknown event")
	}
}

func (r *Runtime) onWiFiFailed() {
	fmt.Println("[EVENT] WiFi failed")

	// Step 1: evaluate decision
	fmt.Println("[DECISION] evaluating alternative path...")
	migrate := true

	if migrate {
		fmt.Println("[DECISION] migrate=true (better path detected)")

		// Step 2: state transition
		r.transition(StateRecovering)

		// Step 3: authority transfer
		r.epoch++
		fmt.Printf("[AUTHORITY] epoch %d granted to 5G\n", r.epoch)

		// Step 4: reject stale transport
		fmt.Println("[CHECK] stale WiFi transport rejected")

		// Step 5: return to attached
		r.transition(StateAttached)

		fmt.Println("[RESULT] session continues (no reconnect, no reset)")
	} else {
		fmt.Println("[RESULT] stay on current transport")
	}
}

func (r *Runtime) transition(newState State) {
	fmt.Printf("[STATE] %s -> %s\n", r.state, newState)
	r.state = newState
}

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO")
	fmt.Println()

	runtime := NewRuntime()

	// simulate failure
	runtime.HandleEvent(EventWiFiFailed)
}