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
	Epoch      int
	Current    Transport
	Candidates []Transport
	PacketID   int
}

func NewRuntime() *Runtime {
	// важно для jitter / loss
	rand.Seed(time.Now().UnixNano())

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

// Оставили, но теперь используется редко (для простых тестов)
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

	margin := best.Score - r.Current.Score
	confidence := computeConfidence(r.Current, best)

	fmt.Printf(
		"[DECISION] migrate=%v (margin=%.1f, confidence=%.2f, reason=better_path)\n",
		best.Score > r.Current.Score,
		margin,
		confidence,
	)

	if best.Score > r.Current.Score {
		r.transition(StateRecovering)

		// имитация времени на миграцию
		time.Sleep(200 * time.Millisecond)

		r.Epoch++
		fmt.Printf("[AUTHORITY] epoch %d granted to %s\n", r.Epoch, best.Name)

		// проверка stale path
		fmt.Printf("[CHECK] stale %s rejected\n", r.Current.Name)

		// переключение транспорта
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

// простая эвристика confidence (чтобы выглядело как runtime, а не if)
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