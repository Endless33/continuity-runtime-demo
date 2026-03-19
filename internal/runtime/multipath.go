package runtime

import (
	"fmt"
)

type MultiPath struct {
	Primary   Transport
	Secondary Transport
	Active    bool
}

func NewMultiPath(primary, secondary Transport) *MultiPath {
	return &MultiPath{
		Primary:   primary,
		Secondary: secondary,
		Active:    false,
	}
}

func (m *MultiPath) StartOverlap() {
	fmt.Println("[MULTIPATH] starting overlap (wifi + 5G)")
	m.Active = true
}

func (m *MultiPath) StopOverlap() {
	fmt.Println("[MULTIPATH] stopping overlap (single path active)")
	m.Active = false
}