package protocol

import "fmt"

type SessionState string

const (
	SessionStateInit       SessionState = "INIT"
	SessionStateAttached   SessionState = "ATTACHED"
	SessionStateRecovering SessionState = "RECOVERING"
	SessionStateClosed     SessionState = "CLOSED"
)

type SessionProtocol struct {
	SessionID      string
	Epoch          int
	State          SessionState
	ActivePath     string
	NextSeq        int
	LastAck        int
	AuthorityOwner string
}

func NewSessionProtocol(sessionID, initialPath string) *SessionProtocol {
	return &SessionProtocol{
		SessionID:      sessionID,
		Epoch:          1,
		State:          SessionStateAttached,
		ActivePath:     initialPath,
		NextSeq:        1,
		LastAck:        0,
		AuthorityOwner: initialPath,
	}
}

func (sp *SessionProtocol) BuildInit() WirePacket {
	return WirePacket{
		Type:      PacketTypeSessionInit,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Path:      sp.ActivePath,
	}
}

func (sp *SessionProtocol) BuildData(payload []byte) WirePacket {
	p := WirePacket{
		Type:      PacketTypeData,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Seq:       sp.NextSeq,
		Path:      sp.ActivePath,
		Payload:   payload,
	}
	sp.NextSeq++
	return p
}

func (sp *SessionProtocol) BuildAck(seq int) WirePacket {
	if seq > sp.LastAck {
		sp.LastAck = seq
	}

	return WirePacket{
		Type:      PacketTypeAck,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Ack:       seq,
		Path:      sp.ActivePath,
	}
}

func (sp *SessionProtocol) StartRecovery(candidatePath string) WirePacket {
	sp.State = SessionStateRecovering

	return WirePacket{
		Type:      PacketTypeAuthorityTransfer,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch + 1,
		Path:      candidatePath,
	}
}

func (sp *SessionProtocol) ApplyAuthorityTransfer(candidatePath string, newEpoch int) error {
	if newEpoch <= sp.Epoch {
		return fmt.Errorf("stale epoch: got %d current %d", newEpoch, sp.Epoch)
	}

	sp.Epoch = newEpoch
	sp.ActivePath = candidatePath
	sp.AuthorityOwner = candidatePath
	sp.State = SessionStateAttached
	return nil
}

func (sp *SessionProtocol) ValidatePacket(pkt WirePacket) error {
	if pkt.SessionID != sp.SessionID {
		return fmt.Errorf("session mismatch")
	}

	if pkt.Epoch != sp.Epoch {
		return fmt.Errorf("stale epoch")
	}

	if pkt.Path != "" && pkt.Path != sp.ActivePath {
		return fmt.Errorf("inactive path")
	}

	return nil
}