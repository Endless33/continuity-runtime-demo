package transport

import (
	"encoding/json"
	"fmt"
)

type PacketType string

const (
	PacketTypeInit    PacketType = "init"
	PacketTypeInitAck PacketType = "init_ack"
	PacketTypeData    PacketType = "data"
)

type Packet struct {
	Type      PacketType `json:"type"`
	SessionID string     `json:"session_id"`
	Epoch     int        `json:"epoch"`
	Payload   string     `json:"payload,omitempty"`
}

func (p Packet) Validate() error {
	if p.Type == "" {
		return fmt.Errorf("missing packet type")
	}
	if p.SessionID == "" {
		return fmt.Errorf("missing session_id")
	}
	if p.Epoch <= 0 {
		return fmt.Errorf("invalid epoch")
	}
	return nil
}

func EncodePacket(p Packet) ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func DecodePacket(b []byte) (Packet, error) {
	var p Packet
	if err := json.Unmarshal(b, &p); err != nil {
		return Packet{}, err
	}
	if err := p.Validate(); err != nil {
		return Packet{}, err
	}
	return p, nil
}