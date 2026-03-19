package runtime

import (
	"fmt"
	"math/rand"
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
	Current    Transport
	Candidates []Transport
	PacketID   int

	Session *Session
	Trace   *TraceRecorder
}

func NewRuntime() *Runtime {
	rand.Seed(time.Now().UnixNano())

	current := Transport{
		Name:    "wifi",
		Latency: 120 * time.Millisecond,
		Score:   20,
	}

	trace := NewTraceRecorder("sess-001")

	rt := &Runtime{
		State: StateAttached,
		Current: current,
		Candidates: []Transport{
			{Name: "5g", Latency: 40 * time.Millisecond, Score: 100},
			{Name: "lte", Latency: 80 * time.Millisecond, Score: 60},
		},
		PacketID: 100,
		Session:  NewSession("sess-001", current),
		Trace:    trace,
	}

	rt.Trace.Record("session_started", "initial transport attached", map[string]interface{}{
		"transport": current.Name,
		"epoch":     rt.Session.Epoch,
	})

	return rt
}

func (r *Runtime) HandleEvent(e Event) {
	switch e {
	case EventWiFiFailed:
		r.onWiFiFailed()
	}
}

func (r *Runtime) onWiFiFailed() {
	fmt.Println("\n[EVENT] WiFi failed")

	r.Trace.Record("transport_failed", "wifi failed", nil)

	best := SelectBestTransport(r.Current, r.Candidates)

	margin := best.Score - r.Current.Score
	confidence := computeConfidence(r.Current, best)

	r.Trace.Record("decision", "migration evaluated", map[string]interface{}{
		"margin":      margin,
		"confidence":  confidence,
		"target_path": best.Name,
	})

	fmt.Printf(
		"[DECISION] migrate=%v (margin=%.1f, confidence=%.2f)\n",
		best.Score > r.Current.Score,
		margin,
		confidence,
	)

	if best.Score > r.Current.Score {
		r.transition(StateRecovering)

		time.Sleep(200 * time.Millisecond)

		r.Session.TransferAuthority(best)

		r.Trace.Record("authority_granted", "authority moved", map[string]interface{}{
			"epoch":     r.Session.Epoch,
			"transport": best.Name,
		})

		fmt.Printf("[AUTHORITY] epoch %d granted to %s\n", r.Session.Epoch, best.Name)

		if !r.Session.ValidateTransport(r.Current.Name, r.Session.Epoch) {
			fmt.Printf("[CHECK] stale %s rejected\n", r.Current.Name)

			r.Trace.Record("stale_rejected", "old path rejected", map[string]interface{}{
				"transport": r.Current.Name,
			})
		}

		r.Current = best

		r.transition(StateAttached)

		r.Trace.Record("session_continued", "continuity preserved", nil)

		fmt.Println("[RESULT] session continues")
	} else {
		fmt.Println("[RESULT] no better transport")
	}
}

func (r *Runtime) transition(newState State) {
	fmt.Printf("[STATE] %s -> %s\n", r.State, newState)

	r.Trace.Record("state_changed", "state transition", map[string]interface{}{
		"from": string(r.State),
		"to":   string(newState),
	})

	r.State = newState
}

func computeConfidence(current Transport, candidate Transport) float64 {
	conf := 0.5

	if candidate.Score > current.Score {
		conf += 0.2
	}
	if candidate.Latency < current.Latency {
		conf += 0.2
	}
	if candidate.Score-current.Score > 50 {
		conf += 0.1
	}

	if conf > 1.0 {
		conf = 1.0
	}

	return conf
}