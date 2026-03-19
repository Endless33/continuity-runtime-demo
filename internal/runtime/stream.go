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
	ACK      *ACKTracker
	RTX      *RetransmitQueue
	Order    *OrderingPolicy
}

func NewStream(r *Runtime, n *NetworkSimulator) *Stream {
	return &Stream{
		Runtime:  r,
		Network:  n,
		Buffer:   NewReorderBuffer(r.PacketID),
		Dedup:    NewDedup(),
		Multi:    NewMultiPath(r.Current, r.Candidates[0]),
		Adaptive: NewAdaptiveController(),
		ACK:      NewACKTracker(),
		RTX:      NewRetransmitQueue(),
		Order:    NewOrderingPolicy(),
	}
}

func (s *Stream) Send(n int) {
	for i := 0; i < n; i++ {
		s.Runtime.PacketID++
		packetID := s.Runtime.PacketID

		if s.Adaptive.Evaluate(s.Runtime.Current) {
			s.Multi.StartOverlap()
		} else {
			s.Multi.StopOverlap()
		}

		if s.Multi.Active {
			winner, ok := RaceTransports(
				s.Network,
				packetID,
				s.Multi.Primary,
				s.Multi.Secondary,
			)
			if !ok {
				s.RTX.Add(Packet{ID: packetID, Transport: s.Runtime.Current.Name})
				continue
			}
			s.process(packetID, winner)
		} else {
			ok := s.Network.Transmit(packetID, s.Runtime.Current)
			if !ok {
				s.RTX.Add(Packet{ID: packetID, Transport: s.Runtime.Current.Name})
				continue
			}
			s.process(packetID, s.Runtime.Current)
		}
	}

	s.retryPending()
}

func (s *Stream) retryPending() {
	s.RTX.Replay(func(p Packet) {
		if s.Dedup.Seen(p.ID) {
			fmt.Printf("[DEDUP] duplicate #%d dropped during retransmit\n", p.ID)
			s.RTX.Remove(p.ID)
			return
		}

		if !s.Order.AllowOutOfOrder(s.Buffer.expectedID, p.ID) {
			fmt.Printf("[ORDER] packet #%d held (expected=%d)\n", p.ID, s.Buffer.expectedID)
			return
		}

		s.Buffer.Push(Packet{
			ID:        p.ID,
			Transport: p.Transport,
		})
		s.ACK.Ack(p.ID)
		s.RTX.Remove(p.ID)
		s.Buffer.DebugState()
	})
}

func (s *Stream) process(packetID int, t Transport) {
	if s.Dedup.Seen(packetID) {
		fmt.Printf("[DEDUP] duplicate #%d dropped\n", packetID)
		return
	}

	if !s.Order.AllowOutOfOrder(s.Buffer.expectedID, packetID) {
		fmt.Printf("[ORDER] packet #%d buffered outside strict order window\n", packetID)
		s.RTX.Add(Packet{ID: packetID, Transport: t.Name})
		return
	}

	s.Buffer.Push(Packet{
		ID:        packetID,
		Transport: t.Name,
	})

	s.ACK.Ack(packetID)
	s.Buffer.DebugState()
}