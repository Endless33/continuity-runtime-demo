package runtime

import "time"

type Transport struct {
	Name    string
	Latency time.Duration
	Score   float64
}