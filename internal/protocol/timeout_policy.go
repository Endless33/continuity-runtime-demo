package protocol

import "time"

type TimeoutPolicy struct {
	InitTimeout         time.Duration
	AckTimeout          time.Duration
	MigrationTimeout    time.Duration
	KeepaliveInterval   time.Duration
	IdleSessionTimeout  time.Duration
}

func NewTimeoutPolicy() *TimeoutPolicy {
	return &TimeoutPolicy{
		InitTimeout:        2 * time.Second,
		AckTimeout:         800 * time.Millisecond,
		MigrationTimeout:   1500 * time.Millisecond,
		KeepaliveInterval:  5 * time.Second,
		IdleSessionTimeout: 30 * time.Second,
	}
}