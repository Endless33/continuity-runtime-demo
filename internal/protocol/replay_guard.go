package protocol

import "fmt"

type ReplayGuard struct {
	highestSeq int
	seen       map[int]bool
	window     int
}

func NewReplayGuard(window int) *ReplayGuard {
	if window <= 0 {
		window = 64
	}

	return &ReplayGuard{
		highestSeq: 0,
		seen:       make(map[int]bool),
		window:     window,
	}
}

func (rg *ReplayGuard) Validate(seq int) error {
	if seq <= 0 {
		return fmt.Errorf("invalid sequence")
	}

	if seq > rg.highestSeq {
		rg.highestSeq = seq
		rg.seen[seq] = true
		rg.gc()
		return nil
	}

	if rg.highestSeq-seq > rg.window {
		return fmt.Errorf("sequence outside replay window")
	}

	if rg.seen[seq] {
		return fmt.Errorf("replay detected")
	}

	rg.seen[seq] = true
	return nil
}

func (rg *ReplayGuard) gc() {
	threshold := rg.highestSeq - rg.window
	for seq := range rg.seen {
		if seq < threshold {
			delete(rg.seen, seq)
		}
	}
}