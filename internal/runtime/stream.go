package runtime

import (
	"fmt"
)

type Stream struct {
	Runtime  *Runtime
	Network  *NetworkSimulator
	Buffer   *ReorderBuffer
	Dedup    *Dedup
	Multi    *MultiPath
	Adaptive *AdaptiveController
}

func NewStream(r *Runtime, n *NetworkSimulator) *Stream {
	return &Stream{
		Runtime:  r,
		Network:  n,
		Buffer:   NewReorderBuffer(r.PacketID),
		Dedup:    NewDedup(),
		Multi:    NewMultiPath(r.Current, r.Candidates[0]),
		Adaptive: NewAdaptiveController(),
	}
}

func (s *Stream) Send(n int) {
	for i := 0; i < n; i++ {
		s.Runtime.PacketID++
		packetID := s.Runtime.PacketID

		// adaptive overlap
		if s.Adaptive.Evaluate(s.Runtime.Current) {
			s.Multi.StartOverlap()
		} else {
			s.Multi.StopOverlap()
		}

		// latency racing если overlap активен
		if s.Multi.Active {
			winner, ok := RaceTransports(
				s.Network,
				packetID,
				s.Multi.Primary,
				s.Multi.Secondary,
			)

			if !ok {
				continue
			}

			s.process(packetID, winner)
		} else {
			ok := s.Network.Transmit(packetID, s.Runtime.Current)
			if !ok {
				continue
			}
			s.process(packetID, s.Runtime.Current)
		}
	}
}

func (s *Stream) process(packetID int, t Transport) {
	// dedup
	if s.Dedup.Seen(packetID) {
		fmt.Printf("[DEDUP] duplicate #%d dropped\n", packetID)
		return
	}

	s.Buffer.Push(Packet{
		ID:        packetID,
		Transport: t.Name,
	})

	s.Buffer.DebugState()
}