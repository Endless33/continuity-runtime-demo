package runtime

import (
	"fmt"
	"math/rand"
)

type Stream struct {
	Runtime *Runtime
	Network *NetworkSimulator
	Buffer  *ReorderBuffer
}

func NewStream(r *Runtime, n *NetworkSimulator) *Stream {
	return &Stream{
		Runtime: r,
		Network: n,
		Buffer:  NewReorderBuffer(r.PacketID),
	}
}

func (s *Stream) Send(n int) {
	for i := 0; i < n; i++ {
		s.Runtime.PacketID++
		packetID := s.Runtime.PacketID

		delivered := s.Network.Transmit(packetID, s.Runtime.Current)

		if !delivered {
			continue
		}

		// имитация out-of-order доставки
		if rand.Float64() < 0.3 {
			delayedID := packetID + rand.Intn(3)

			fmt.Printf("[REORDER] packet #%d delayed → arrives as #%d\n", packetID, delayedID)

			s.Buffer.Push(Packet{
				ID:        delayedID,
				Transport: s.Runtime.Current.Name,
			})
		} else {
			s.Buffer.Push(Packet{
				ID:        packetID,
				Transport: s.Runtime.Current.Name,
			})
		}

		s.Buffer.DebugState()
	}
}