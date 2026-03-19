package runtime

import (
	"fmt"
	"sort"
)

type Packet struct {
	ID        int
	Transport string
}

type ReorderBuffer struct {
	expectedID int
	buffer     map[int]Packet
}

func NewReorderBuffer(startID int) *ReorderBuffer {
	return &ReorderBuffer{
		expectedID: startID + 1,
		buffer:     make(map[int]Packet),
	}
}

func (rb *ReorderBuffer) Push(p Packet) {
	rb.buffer[p.ID] = p
	rb.flush()
}

func (rb *ReorderBuffer) flush() {
	for {
		p, ok := rb.buffer[rb.expectedID]
		if !ok {
			break
		}

		fmt.Printf("[DELIVER] packet #%d (ordered, via %s)\n", p.ID, p.Transport)

		delete(rb.buffer, rb.expectedID)
		rb.expectedID++
	}
}

func (rb *ReorderBuffer) DebugState() {
	if len(rb.buffer) == 0 {
		return
	}

	keys := make([]int, 0, len(rb.buffer))
	for k := range rb.buffer {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	fmt.Printf("[BUFFER] waiting for #%d, have: %v\n", rb.expectedID, keys)
}