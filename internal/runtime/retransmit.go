package runtime

import "fmt"

type RetransmitQueue struct {
	pending map[int]Packet
}

func NewRetransmitQueue() *RetransmitQueue {
	return &RetransmitQueue{
		pending: make(map[int]Packet),
	}
}

func (rq *RetransmitQueue) Add(p Packet) {
	rq.pending[p.ID] = p
	fmt.Printf("[RTX] queued packet #%d for retransmission\n", p.ID)
}

func (rq *RetransmitQueue) Remove(id int) {
	delete(rq.pending, id)
}

func (rq *RetransmitQueue) Replay(process func(Packet)) {
	for _, p := range rq.pending {
		fmt.Printf("[RTX] retransmitting packet #%d via %s\n", p.ID, p.Transport)
		process(p)
	}
}