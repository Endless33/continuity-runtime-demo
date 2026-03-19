package protocol

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

	Replay *ReplayGuard
	Window *SeqWindow
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
		Replay:         NewReplayGuard(64),
		Window:         NewSeqWindow(),
	}
}

func (sp *SessionProtocol) BuildInit() WirePacket {
	return WirePacket{
		Version:   ProtocolVersion,
		Type:      PacketTypeSessionInit,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch,
		Path:      sp.ActivePath,
	}
}

func (sp *SessionProtocol) BuildData(payload []byte) WirePacket {
	p := WirePacket{
		Version:   ProtocolVersion,
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
		Version:   ProtocolVersion,
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
		Version:   ProtocolVersion,
		Type:      PacketTypeAuthorityTransfer,
		SessionID: sp.SessionID,
		Epoch:     sp.Epoch + 1,
		Path:      candidatePath,
	}
}

func (sp *SessionProtocol) ApplyAuthorityTransfer(candidatePath string, newEpoch int) error {
	if sp.State == SessionStateClosed {
		return NewProtocolError(ErrProtocolClosed, "cannot migrate closed session")
	}

	if newEpoch <= sp.Epoch {
		return NewProtocolError(ErrStaleEpoch, "authority transfer epoch is not newer")
	}

	sp.Epoch = newEpoch
	sp.ActivePath = candidatePath
	sp.AuthorityOwner = candidatePath
	sp.State = SessionStateAttached
	return nil
}

func (sp *SessionProtocol) ValidatePacket(pkt WirePacket) error {
	if sp.State == SessionStateClosed && pkt.Type != PacketTypeClose {
		return NewProtocolError(ErrProtocolClosed, "session is closed")
	}

	if pkt.Version != ProtocolVersion {
		return NewProtocolError(ErrUnsupportedVersion, "unsupported protocol version")
	}

	if pkt.SessionID != sp.SessionID {
		return NewProtocolError(ErrSessionMismatch, "session mismatch")
	}

	if pkt.Epoch != sp.Epoch {
		return NewProtocolError(ErrStaleEpoch, "stale epoch")
	}

	if pkt.Path != "" && pkt.Path != sp.ActivePath {
		return NewProtocolError(ErrInactivePath, "inactive path")
	}

	switch pkt.Type {
	case PacketTypeData:
		if err := sp.Window.Validate(pkt.Seq); err != nil {
			return NewProtocolError(ErrInvalidSequence, err.Error())
		}
		if err := sp.Replay.Validate(pkt.Seq); err != nil {
			return NewProtocolError(ErrReplayDetected, err.Error())
		}
		sp.Window.Advance(pkt.Seq)

	case PacketTypeAck:
		if pkt.Ack <= 0 {
			return NewProtocolError(ErrAckInvalid, "invalid ack")
		}

	case PacketTypeKeepalive:
		return nil

	case PacketTypeClose:
		sp.State = SessionStateClosed
		return nil
	}

	return nil
}