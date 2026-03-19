package protocol

import "encoding/json"

type PacketType string

const (
	PacketTypeData              PacketType = "data"
	PacketTypeAck               PacketType = "ack"
	PacketTypeSessionInit       PacketType = "session_init"
	PacketTypeSessionInitAck    PacketType = "session_init_ack"
	PacketTypeAuthorityTransfer PacketType = "authority_transfer"
	PacketTypeAuthorityAck      PacketType = "authority_ack"
	PacketTypeKeepalive         PacketType = "keepalive"
)

type WirePacket struct {
	Type      PacketType `json:"type"`
	SessionID string     `json:"session_id"`
	Epoch     int        `json:"epoch"`
	Seq       int        `json:"seq"`
	Ack       int        `json:"ack,omitempty"`
	Path      string     `json:"path,omitempty"`
	Payload   []byte     `json:"payload,omitempty"`
}

func EncodePacket(p WirePacket) ([]byte, error) {
	return json.Marshal(p)
}

func DecodePacket(b []byte) (WirePacket, error) {
	var p WirePacket
	err := json.Unmarshal(b, &p)
	return p, err
}