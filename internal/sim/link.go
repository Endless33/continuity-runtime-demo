package sim

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
)

type Link struct {
	Name string
}

func NewLink(name string) *Link {
	return &Link{Name: name}
}

func (l *Link) Deliver(pkt protocol.WirePacket, receive func(protocol.WirePacket) error) error {
	fmt.Printf("[LINK %s] deliver type=%s seq=%d epoch=%d path=%s\n",
		l.Name,
		pkt.Type,
		pkt.Seq,
		pkt.Epoch,
		pkt.Path,
	)

	return receive(pkt)
}