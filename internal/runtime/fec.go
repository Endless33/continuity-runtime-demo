package runtime

import "fmt"

type FECBlock struct {
	BlockID      int
	DataPackets  []int
	ParityPacket int
}

type FECEncoder struct {
	nextBlockID int
	groupSize   int
}

func NewFECEncoder(groupSize int) *FECEncoder {
	return &FECEncoder{
		nextBlockID: 1,
		groupSize:   groupSize,
	}
}

func (f *FECEncoder) Build(packets []Packet) *FECBlock {
	if len(packets) < f.groupSize {
		return nil
	}

	block := &FECBlock{
		BlockID:     f.nextBlockID,
		DataPackets: make([]int, 0, f.groupSize),
	}

	parity := 0
	for _, p := range packets[:f.groupSize] {
		block.DataPackets = append(block.DataPackets, p.ID)
		parity ^= p.ID
	}

	block.ParityPacket = parity
	f.nextBlockID++

	fmt.Printf("[FEC] block #%d built data=%v parity=%d\n",
		block.BlockID,
		block.DataPackets,
		block.ParityPacket,
	)

	return block
}

type FECDecoder struct{}

func NewFECDecoder() *FECDecoder {
	return &FECDecoder{}
}

func (d *FECDecoder) Recover(block FECBlock, received []int) (int, bool) {
	if len(received) != len(block.DataPackets)-1 {
		return 0, false
	}

	value := block.ParityPacket
	for _, id := range received {
		value ^= id
	}

	fmt.Printf("[FEC] recovered missing packet #%d from block #%d\n", value, block.BlockID)
	return value, true
}