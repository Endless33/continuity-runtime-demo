package sim

import (
	"fmt"
	"math/rand"
	"time"

	"continuity-runtime-demo/internal/protocol"
)

type LossyLink struct {
	Name      string
	LossRate  float64
	MaxDelay  time.Duration
	Reorder   bool
	Duplicate bool
}

func NewLossyLink(name string) *LossyLink {
	rand.Seed(time.Now().UnixNano())

	return &LossyLink{
		Name:      name,
		LossRate:  0.20,
		MaxDelay:  120 * time.Millisecond,
		Reorder:   true,
		Duplicate: true,
	}
}

func (l *LossyLink) Deliver(pkt protocol.WirePacket, receive func(protocol.WirePacket) error) error {
	if rand.Float64() < l.LossRate {
		fmt.Printf("[LOSSY %s] drop type=%s seq=%d epoch=%d path=%s\n",
			l.Name, pkt.Type, pkt.Seq, pkt.Epoch, pkt.Path,
		)
		return nil
	}

	delay := time.Duration(rand.Int63n(int64(l.MaxDelay) + 1))
	if delay > 0 {
		time.Sleep(delay)
	}

	fmt.Printf("[LOSSY %s] deliver type=%s seq=%d epoch=%d path=%s delay=%v\n",
		l.Name, pkt.Type, pkt.Seq, pkt.Epoch, pkt.Path, delay,
	)

	if err := receive(pkt); err != nil {
		return err
	}

	if l.Duplicate && rand.Float64() < 0.15 {
		fmt.Printf("[LOSSY %s] duplicate type=%s seq=%d epoch=%d path=%s\n",
			l.Name, pkt.Type, pkt.Seq, pkt.Epoch, pkt.Path,
		)
		if err := receive(pkt); err != nil {
			return err
		}
	}

	if l.Reorder && rand.Float64() < 0.15 {
		extraDelay := time.Duration(rand.Int63n(int64(l.MaxDelay) + 1))
		time.Sleep(extraDelay)
		fmt.Printf("[LOSSY %s] reordered replay type=%s seq=%d epoch=%d path=%s delay=%v\n",
			l.Name, pkt.Type, pkt.Seq, pkt.Epoch, pkt.Path, extraDelay,
		)
		if err := receive(pkt); err != nil {
			return err
		}
	}

	return nil
}