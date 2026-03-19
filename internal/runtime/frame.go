package runtime

import "fmt"

type Frame struct {
	FrameID    int
	PacketIDs  []int
	Complete   bool
	Recovered  bool
	Transport  string
}

type FrameAssembler struct {
	nextFrameID int
	frameSize   int
	pending     []Packet
}

func NewFrameAssembler(frameSize int) *FrameAssembler {
	return &FrameAssembler{
		nextFrameID: 1,
		frameSize:   frameSize,
		pending:     []Packet{},
	}
}

func (fa *FrameAssembler) Push(p Packet) *Frame {
	fa.pending = append(fa.pending, p)

	if len(fa.pending) < fa.frameSize {
		return nil
	}

	frame := &Frame{
		FrameID:   fa.nextFrameID,
		PacketIDs: make([]int, 0, fa.frameSize),
		Complete:  true,
		Recovered: false,
		Transport: p.Transport,
	}

	for _, pkt := range fa.pending[:fa.frameSize] {
		frame.PacketIDs = append(frame.PacketIDs, pkt.ID)
	}

	fa.pending = fa.pending[fa.frameSize:]
	fa.nextFrameID++

	fmt.Printf("[FRAME] frame #%d assembled from packets %v\n", frame.FrameID, frame.PacketIDs)

	return frame
}

func (fa *FrameAssembler) FlushPartial() *Frame {
	if len(fa.pending) == 0 {
		return nil
	}

	frame := &Frame{
		FrameID:   fa.nextFrameID,
		PacketIDs: make([]int, 0, len(fa.pending)),
		Complete:  false,
		Recovered: false,
		Transport: fa.pending[0].Transport,
	}

	for _, pkt := range fa.pending {
		frame.PacketIDs = append(frame.PacketIDs, pkt.ID)
	}

	fa.pending = nil
	fa.nextFrameID++

	fmt.Printf("[FRAME] partial frame #%d flushed from packets %v\n", frame.FrameID, frame.PacketIDs)

	return frame
}