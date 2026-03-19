package runtime

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
)

type Exchange struct {
	Name string
}

func NewExchange(name string) *Exchange {
	return &Exchange{
		Name: name,
	}
}

func (ex *Exchange) Send(from, to *Node, pkt protocol.WirePacket) (*protocol.WirePacket, error) {
	if from == nil || to == nil {
		return nil, fmt.Errorf("nil node")
	}

	from.Engine.Runtime.Trace.Record("exchange_send", "packet sent", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"type":     pkt.Type,
		"seq":      pkt.Seq,
		"epoch":    pkt.Epoch,
		"path":     pkt.Path,
	})

	fmt.Printf("[EXCHANGE %s] %s -> %s type=%s seq=%d epoch=%d path=%s\n",
		ex.Name,
		from.Name,
		to.Name,
		pkt.Type,
		pkt.Seq,
		pkt.Epoch,
		pkt.Path,
	)

	resp, err := to.Receive(pkt)
	if err != nil {
		to.Engine.Runtime.Trace.Record("exchange_receive_error", "packet rejected", map[string]interface{}{
			"exchange": ex.Name,
			"from":     from.Name,
			"to":       to.Name,
			"type":     pkt.Type,
			"seq":      pkt.Seq,
			"epoch":    pkt.Epoch,
			"path":     pkt.Path,
			"error":    err.Error(),
		})
		return nil, err
	}

	to.Engine.Runtime.Trace.Record("exchange_receive", "packet received", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"type":     pkt.Type,
		"seq":      pkt.Seq,
		"epoch":    pkt.Epoch,
		"path":     pkt.Path,
	})

	return resp, nil
}

func (ex *Exchange) RoundTrip(a, b *Node, pkt protocol.WirePacket) error {
	resp, err := ex.Send(a, b, pkt)
	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}

	_, err = ex.Send(b, a, *resp)
	return err
}

func (ex *Exchange) SendAck(from, to *Node, ackFor int) error {
	if from == nil || to == nil {
		return fmt.Errorf("nil node")
	}

	ack := from.Engine.Protocol.BuildAck(ackFor)

	from.Engine.Runtime.Trace.Record("exchange_ack_sent", "ack sent", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"ack":      ack.Ack,
		"epoch":    ack.Epoch,
		"path":     ack.Path,
	})

	fmt.Printf("[EXCHANGE %s] %s -> %s ACK=%d epoch=%d path=%s\n",
		ex.Name,
		from.Name,
		to.Name,
		ack.Ack,
		ack.Epoch,
		ack.Path,
	)

	if err := to.Engine.Receive(ack); err != nil {
		to.Engine.Runtime.Trace.Record("exchange_ack_error", "ack rejected", map[string]interface{}{
			"exchange": ex.Name,
			"from":     from.Name,
			"to":       to.Name,
			"ack":      ack.Ack,
			"epoch":    ack.Epoch,
			"path":     ack.Path,
			"error":    err.Error(),
		})
		return err
	}

	to.Engine.Runtime.Trace.Record("exchange_ack_received", "ack received", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"ack":      ack.Ack,
		"epoch":    ack.Epoch,
		"path":     ack.Path,
	})

	return nil
}