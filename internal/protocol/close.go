package protocol

type CloseReason string

const (
	CloseReasonNormal        CloseReason = "normal"
	CloseReasonTimeout       CloseReason = "timeout"
	CloseReasonProtocolError CloseReason = "protocol_error"
	CloseReasonAuthorityLost CloseReason = "authority_lost"
)

type CloseEngine struct{}

func NewCloseEngine() *CloseEngine {
	return &CloseEngine{}
}

func (c *CloseEngine) Build(sp *SessionProtocol, reason CloseReason) WirePacket {
	return WirePacket{
		Version:   ProtocolVersion,
		Type:      PacketTypeClose,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Path:      sp.ActivePath,
		Payload:   []byte(reason),
	}
}

func (c *CloseEngine) Apply(sp *SessionProtocol) {
	if sp == nil {
		return
	}
	sp.State = SessionStateClosed
}

func DecodeCloseReason(pkt WirePacket) CloseReason {
	if len(pkt.Payload) == 0 {
		return CloseReasonNormal
	}
	return CloseReason(pkt.Payload)
}