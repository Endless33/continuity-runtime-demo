package protocol

import "time"

type KeepaliveEngine struct {
	Interval time.Duration
	LastSeen time.Time
	LastSent time.Time
}

func NewKeepaliveEngine(interval time.Duration) *KeepaliveEngine {
	if interval <= 0 {
		interval = 5 * time.Second
	}

	now := time.Now().UTC()

	return &KeepaliveEngine{
		Interval: interval,
		LastSeen: now,
		LastSent: now,
	}
}

func (k *KeepaliveEngine) ShouldSend(now time.Time) bool {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return now.Sub(k.LastSent) >= k.Interval
}

func (k *KeepaliveEngine) MarkSent(now time.Time) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	k.LastSent = now
}

func (k *KeepaliveEngine) MarkSeen(now time.Time) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	k.LastSeen = now
}

func (k *KeepaliveEngine) Build(sp *SessionProtocol) WirePacket {
	return WirePacket{
		Version:   ProtocolVersion,
		Type:      PacketTypeKeepalive,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Path:      sp.ActivePath,
	}
}

func (k *KeepaliveEngine) IsExpired(now time.Time, idleTimeout time.Duration) bool {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if idleTimeout <= 0 {
		return false
	}
	return now.Sub(k.LastSeen) > idleTimeout
}