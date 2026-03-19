package runtime

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
)

type AckExchange struct {
	Forward  *LossyExchange
	Reverse  *LossyExchange
	AckRules *protocol.AckRules
}

func NewAckExchange(forward, reverse *LossyExchange) *AckExchange {
	return &AckExchange{
		Forward:  forward,
		Reverse:  reverse,
		AckRules: protocol.NewAckRules(),
	}
}

func (ax *AckExchange) SendData(from, to *Node, pkt protocol.WirePacket) error {
	if from == nil || to == nil {
		return fmt.Errorf("nil node")
	}

	if pkt.Type != protocol.PacketTypeData {
		return fmt.Errorf("packet is not data")
	}

	_, err := ax.Forward.Send(from, to, pkt)
	if err != nil {
		return err
	}

	ack := to.Engine.Protocol.BuildAck(pkt.Seq)

	_, err = ax.Reverse.Send(to, from, ack)
	if err != nil {
		return err
	}

	res := ax.AckRules.ValidateAck(from.Engine.Protocol, ack)
	if res.Allowed {
		fmt.Printf("[ACK RULES] accepted ack=%d reason=%s\n", ack.Ack, res.Reason)
	} else {
		fmt.Printf("[ACK RULES] rejected ack=%d reason=%s\n", ack.Ack, res.Reason)
	}

	from.Engine.Runtime.Trace.Record("ack_exchange_result", "ack evaluated", map[string]interface{}{
		"ack":      ack.Ack,
		"allowed":  res.Allowed,
		"decision": res.Decision,
		"reason":   res.Reason,
	})

	return nil
}

func (ax *AckExchange) SendDataBatch(from, to *Node, packets []protocol.WirePacket) error {
	for _, pkt := range packets {
		if err := ax.SendData(from, to, pkt); err != nil {
			return err
		}
	}
	return nil
}