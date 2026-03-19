package runtime

import (
	"fmt"
	"math/rand"
)

type Stream struct {
	Runtime *Runtime
	Network *NetworkSimulator
	Buffer  *ReorderBuffer
	Dedup   *Dedup
	Multi   *MultiPath
}

func NewStream(r *Runtime, n *NetworkSimulator) *Stream {
	return &Stream{
		Runtime: r,
		Network: n,
		Buffer:  NewReorderBuffer(r.PacketID),
		Dedup:   NewDedup(),
		Multi:   NewMultiPath(r.Current, r.Candidates[0]),
	}
}

func (s *Stream) Send(n int) {
	for i := 0; i < n; i++ {
		s.Runtime.PacketID++
		packetID := s.Runtime.PacketID

		// если overlap активен — отправляем по двум путям
		if s.Multi.Active {
			s.sendVia(packetID, s.Multi.Primary)
			s.sendVia(packetID, s.Multi.Secondary)
		} else {
			s.sendVia(packetID, s.Runtime.Current)
		}
	}
}

func (s *Stream) sendVia(packetID int, t Transport) {
	delivered := s.Network.Transmit(packetID, t)

	if !delivered {
		return
	}

	// dedup (важно!)
	if s.Dedup.Seen(packetID) {
		fmt.Printf("[DEDUP] duplicate packet #%d dropped\n", packetID)
		return
	}

	// reorder simulation
	if rand.Float64() < 0.3 {
		delayedID := packetID + rand.Intn(3)

		fmt.Printf("[REORDER] packet #%d delayed → arrives as #%d (%s)\n",
			packetID, delayedID, t.Name)

		s.Buffer.Push(Packet{
			ID:        delayedID,
			Transport: t.Name,
		})
	} else {
		s.Buffer.Push(Packet{
			ID:        packetID,
			Transport: t.Name,
		})
	}

	s.Buffer.DebugState()
}