package protocol

import "fmt"

type AckDecision string

const (
	AckAccept AckDecision = "accept_ack"
	AckReject AckDecision = "reject_ack"
)

type AckResult struct {
	Decision AckDecision
	Allowed  bool
	Reason   string
}

type AckRules struct{}

func NewAckRules() *AckRules {
	return &AckRules{}
}

func (ar *AckRules) ValidateAck(sp *SessionProtocol, pkt WirePacket) AckResult {
	if sp == nil {
		return AckResult{
			Decision: AckReject,
			Allowed:  false,
			Reason:   "nil session protocol",
		}
	}

	if pkt.Type != PacketTypeAck {
		return AckResult{
			Decision: AckReject,
			Allowed:  false,
			Reason:   "packet is not ack",
		}
	}

	if pkt.SessionID != sp.SessionID {
		return AckResult{
			Decision: AckReject,
			Allowed:  false,
			Reason:   "session mismatch",
		}
	}

	if pkt.Epoch != sp.Epoch {
		return AckResult{
			Decision: AckReject,
			Allowed:  false,
			Reason:   fmt.Sprintf("ack epoch %d != session epoch %d", pkt.Epoch, sp.Epoch),
		}
	}

	if pkt.Ack <= 0 {
		return AckResult{
			Decision: AckReject,
			Allowed:  false,
			Reason:   "invalid ack number",
		}
	}

	if pkt.Ack > sp.NextSeq-1 {
		return AckResult{
			Decision: AckReject,
			Allowed:  false,
			Reason:   fmt.Sprintf("ack %d exceeds highest sent seq %d", pkt.Ack, sp.NextSeq-1),
		}
	}

	return AckResult{
		Decision: AckAccept,
		Allowed:  true,
		Reason:   "ack accepted",
	}
}