package protocol

import "fmt"

type SeqWindow struct {
	Expected int
	Window   int
}

func NewSeqWindow() *SeqWindow {
	return &SeqWindow{
		Expected: 1,
		Window:   32,
	}
}

func (w *SeqWindow) Validate(seq int) error {
	if seq <= 0 {
		return fmt.Errorf("invalid sequence")
	}

	if seq < w.Expected-w.Window {
		return fmt.Errorf("sequence too old")
	}

	if seq > w.Expected+w.Window {
		return fmt.Errorf("sequence too far ahead")
	}

	return nil
}

func (w *SeqWindow) Advance(seq int) {
	if seq >= w.Expected {
		w.Expected = seq + 1
	}
}