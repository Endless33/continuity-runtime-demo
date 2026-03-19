package runtime

import (
	"fmt"
	"math/rand"
	"time"

	"continuity-runtime-demo/internal/protocol"
)

type LossyExchange struct {
	Name         string
	LossRate     float64
	DuplicateRate float64
	MaxDelay     time.Duration
}

func NewLossyExchange(name string) *LossyExchange {
	rand.Seed(time.Now().UnixNano())

	return &LossyExchange{
		Name:          name,
		LossRate:      0.20,
		DuplicateRate: 0.15,
		MaxDelay:      120 * time.Millisecond,
	}
}

func (ex *LossyExchange) Send(from, to *Node, pkt protocol.WirePacket) (*protocol.WirePacket, error) {
	if from == nil || to == nil {
		return nil, fmt.Errorf("nil node")
	}

	from.Engine.Runtime.Trace.Record("lossy_exchange_send", "packet sent", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"type":     pkt.Type,
		"seq":      pkt.Seq,
		"epoch":    pkt.Epoch,
		"path":     pkt.Path,
	})

	if rand.Float64() < ex.LossRate {
		fmt.Printf("[LOSSY %s] drop %s -> %s type=%s seq=%d epoch=%d path=%s\n",
			ex.Name,
			from.Name,
			to.Name,
			pkt.Type,
			pkt.Seq,
			pkt.Epoch,
			pkt.Path,
		)

		to.Engine.Runtime.Trace.Record("lossy_exchange_drop", "packet dropped", map[string]interface{}{
			"exchange": ex.Name,
			"from":     from.Name,
			"to":       to.Name,
			"type":     pkt.Type,
			"seq":      pkt.Seq,
			"epoch":    pkt.Epoch,
			"path":     pkt.Path,
		})

		return nil, nil
	}

	delay := time.Duration(rand.Int63n(int64(ex.MaxDelay) + 1))
	if delay > 0 {
		time.Sleep(delay)
	}

	fmt.Printf("[LOSSY %s] deliver %s -> %s type=%s seq=%d epoch=%d path=%s delay=%v\n",
		ex.Name,
		from.Name,
		to.Name,
		pkt.Type,
		pkt.Seq,
		pkt.Epoch,
		pkt.Path,
		delay,
	)

	resp, err := to.Receive(pkt)
	if err != nil {
		to.Engine.Runtime.Trace.Record("lossy_exchange_receive_error", "packet rejected", map[string]interface{}{
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

	to.Engine.Runtime.Trace.Record("lossy_exchange_receive", "packet received", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"type":     pkt.Type,
		"seq":      pkt.Seq,
		"epoch":    pkt.Epoch,
		"path":     pkt.Path,
		"delay_ms": delay.Milliseconds(),
	})

	if rand.Float64() < ex.DuplicateRate {
		fmt.Printf("[LOSSY %s] duplicate %s -> %s type=%s seq=%d epoch=%d path=%s\n",
			ex.Name,
			from.Name,
			to.Name,
			pkt.Type,
			pkt.Seq,
			pkt.Epoch,
			pkt.Path,
		)

		_, _ = to.Receive(pkt)

		to.Engine.Runtime.Trace.Record("lossy_exchange_duplicate", "duplicate packet delivered", map[string]interface{}{
			"exchange": ex.Name,
			"from":     from.Name,
			"to":       to.Name,
			"type":     pkt.Type,
			"seq":      pkt.Seq,
			"epoch":    pkt.Epoch,
			"path":     pkt.Path,
		})
	}

	return resp, nil
}

func (ex *LossyExchange) RoundTrip(a, b *Node, pkt protocol.WirePacket) error {
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

func (ex *LossyExchange) SendAck(from, to *Node, ackFor int) error {
	if from == nil || to == nil {
		return fmt.Errorf("nil node")
	}

	ack := from.Engine.Protocol.BuildAck(ackFor)

	from.Engine.Runtime.Trace.Record("lossy_exchange_ack_send", "ack sent", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"ack":      ack.Ack,
		"epoch":    ack.Epoch,
		"path":     ack.Path,
	})

	if rand.Float64() < ex.LossRate {
		fmt.Printf("[LOSSY %s] drop ACK %s -> %s ack=%d epoch=%d path=%s\n",
			ex.Name,
			from.Name,
			to.Name,
			ack.Ack,
			ack.Epoch,
			ack.Path,
		)

		to.Engine.Runtime.Trace.Record("lossy_exchange_ack_drop", "ack dropped", map[string]interface{}{
			"exchange": ex.Name,
			"from":     from.Name,
			"to":       to.Name,
			"ack":      ack.Ack,
			"epoch":    ack.Epoch,
			"path":     ack.Path,
		})

		return nil
	}

	delay := time.Duration(rand.Int63n(int64(ex.MaxDelay) + 1))
	if delay > 0 {
		time.Sleep(delay)
	}

	fmt.Printf("[LOSSY %s] deliver ACK %s -> %s ack=%d epoch=%d path=%s delay=%v\n",
		ex.Name,
		from.Name,
		to.Name,
		ack.Ack,
		ack.Epoch,
		ack.Path,
		delay,
	)

	if err := to.Engine.Receive(ack); err != nil {
		to.Engine.Runtime.Trace.Record("lossy_exchange_ack_error", "ack rejected", map[string]interface{}{
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

	to.Engine.Runtime.Trace.Record("lossy_exchange_ack_received", "ack received", map[string]interface{}{
		"exchange": ex.Name,
		"from":     from.Name,
		"to":       to.Name,
		"ack":      ack.Ack,
		"epoch":    ack.Epoch,
		"path":     ack.Path,
		"delay_ms": delay.Milliseconds(),
	})

	return nil
}