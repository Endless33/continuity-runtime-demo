package protocol

import "fmt"

type HandshakeResult struct {
	Accepted bool
	Reason   string
	Packet    WirePacket
}

type HandshakeEngine struct{}

func NewHandshakeEngine() *HandshakeEngine {
	return &HandshakeEngine{}
}

func (h *HandshakeEngine) BuildInit(sp *SessionProtocol) WirePacket {
	return WirePacket{
		Type:      PacketTypeSessionInit,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Path:      sp.ActivePath,
	}
}

func (h *HandshakeEngine) HandleInit(sp *SessionProtocol, pkt WirePacket) HandshakeResult {
	if pkt.Type != PacketTypeSessionInit {
		return HandshakeResult{
			Accepted: false,
			Reason:   "not a session_init packet",
		}
	}

	if pkt.SessionID == "" {
		return HandshakeResult{
			Accepted: false,
			Reason:   "empty session id",
		}
	}

	sp.SessionID = pkt.SessionID
	sp.Epoch = pkt.Epoch
	sp.ActivePath = pkt.Path
	sp.AuthorityOwner = pkt.Path
	sp.State = SessionStateAttached

	ack := WirePacket{
		Type:      PacketTypeSessionInitAck,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Path:      sp.ActivePath,
	}

	return HandshakeResult{
		Accepted: true,
		Reason:   "session initialized",
		Packet:   ack,
	}
}

func (h *HandshakeEngine) HandleInitAck(sp *SessionProtocol, pkt WirePacket) HandshakeResult {
	if pkt.Type != PacketTypeSessionInitAck {
		return HandshakeResult{
			Accepted: false,
			Reason:   "not a session_init_ack packet",
		}
	}

	if pkt.SessionID != sp.SessionID {
		return HandshakeResult{
			Accepted: false,
			Reason:   "session mismatch",
		}
	}

	if pkt.Epoch != sp.Epoch {
		return HandshakeResult{
			Accepted: false,
			Reason:   fmt.Sprintf("epoch mismatch: got %d expected %d", pkt.Epoch, sp.Epoch),
		}
	}

	sp.State = SessionStateAttached
	sp.ActivePath = pkt.Path
	sp.AuthorityOwner = pkt.Path

	return HandshakeResult{
		Accepted: true,
		Reason:   "session attached",
	}
}