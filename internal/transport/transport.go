package transport

import "time"

type PacketIO interface {
	Send([]byte) error
	Receive() ([]byte, error)
}

type Meta struct {
	Name    string
	Latency time.Duration
	Healthy bool
	Score   float64
}

type Transport interface {
	PacketIO
	Meta() Meta
}