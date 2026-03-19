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

	Frames *FrameAssembler
	FECEnc *FECEncoder
	FECDec *FECDecoder
	Recent []Packet
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

		Frames: NewFrameAssembler(3),
		FECEnc: NewFECEncoder(3),
		FECDec: NewFECDecoder(),
		Recent: []Packet{},
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
				s.handleLoss(packetID, s.Runtime.Current.Name)
				continue
			}
			s.process(packetID, winner)
		} else {
			ok := s.Network.Transmit(packetID, s.Runtime.Current)
			if !ok {
				s.handleLoss(packetID, s.Runtime.Current.Name)
				continue
			}
			s.process(packetID, s.Runtime.Current)
		}
	}

	s.retryPending()

	// в конце можно флашить неполный frame
	if frame := s.Frames.FlushPartial(); frame != nil {
		fmt.Printf("[FRAME] emitted partial frame #%d (complete=%v)\n", frame.FrameID, frame.Complete)
	}
}

func (s *Stream) handleLoss(packetID int, transport string) {
	p := Packet{
		ID:        packetID,
		Transport: transport,
	}

	s.RTX.Add(p)

	// Попытка FEC recovery, если хватает контекста
	if len(s.Recent) >= 3 {
		block := s.FECEnc.Build(s.Recent[:3])
		if block != nil {
			received := []int{}
			missingID := packetID

			for _, id := range block.DataPackets {
				if id != missingID {
					received = append(received, id)
				}
			}

			if recovered, ok := s.FECDec.Recover(*block, received); ok {
				fmt.Printf("[FEC] packet #%d recovered logically\n", recovered)

				s.process(recovered, Transport{
					Name:    transport + "+fec",
					Latency: 0,
					Score:   s.Runtime.Current.Score,
				})

				s.RTX.Remove(packetID)
			}
		}
	}
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
		s.captureForFrame(Packet{
			ID:        p.ID,
			Transport: p.Transport,
		})
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

	p := Packet{
		ID:        packetID,
		Transport: t.Name,
	}

	s.Buffer.Push(p)
	s.ACK.Ack(packetID)
	s.captureForFrame(p)
	s.Buffer.DebugState()
}

func (s *Stream) captureForFrame(p Packet) {
	s.Recent = append(s.Recent, p)
	if len(s.Recent) > 6 {
		s.Recent = s.Recent[len(s.Recent)-6:]
	}

	if frame := s.Frames.Push(p); frame != nil {
		fmt.Printf("[FRAME] ready frame #%d complete=%v recovered=%v\n",
			frame.FrameID,
			frame.Complete,
			frame.Recovered,
		)
	}

	if len(s.Recent) >= 3 {
		_ = s.FECEnc.Build(s.Recent[len(s.Recent)-3:])
	}
}