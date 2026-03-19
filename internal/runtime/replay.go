package runtime

import (
	"fmt"
)

type ReplayEngine struct {
	Events []TraceEvent
}

func NewReplayEngine(events []TraceEvent) *ReplayEngine {
	return &ReplayEngine{
		Events: events,
	}
}

func (r *ReplayEngine) Run() {
	fmt.Println("\n=== REPLAY ===")

	for i, ev := range r.Events {
		fmt.Printf("[REPLAY T+%02d] %s → %s\n", i, ev.Type, ev.Message)
	}
}