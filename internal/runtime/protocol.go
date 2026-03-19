package runtime

type ProtocolEvent struct {
	Type    string
	Session string
	Epoch   int
	Path    string
}

func FromTrace(ev TraceEvent) ProtocolEvent {
	return ProtocolEvent{
		Type:    ev.Type,
		Session: ev.SessionID,
	}
}