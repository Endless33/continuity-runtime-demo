package runtime

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
)

type Engine struct {
	Runtime   *Runtime
	Protocol  *protocol.SessionProtocol
	Invariant *protocol.InvariantChecker
}

func NewEngine() *Engine {
	rt := NewRuntime()

	return &Engine{
		Runtime:   rt,
		Protocol:  protocol.NewSessionProtocol("sess-001", rt.Current.Name),
		Invariant: protocol.NewInvariantChecker(),
	}
}

func (e *Engine) Init() {
	pkt := e.Protocol.BuildInit()

	e.Runtime.Trace.Record("protocol_init", "session protocol initialized", map[string]interface{}{
		"session_id": pkt.SessionID,
		"epoch":      pkt.Epoch,
		"path":       pkt.Path,
	})

	fmt.Printf("[PROTOCOL] init session=%s epoch=%d path=%s\n",
		pkt.SessionID,
		pkt.Epoch,
		pkt.Path,
	)
}

func (e *Engine) SendData(payload []byte) {
	pkt := e.Protocol.BuildData(payload)

	e.Runtime.Trace.Record("protocol_data", "data packet built", map[string]interface{}{
		"seq":   pkt.Seq,
		"epoch": pkt.Epoch,
		"path":  pkt.Path,
	})

	fmt.Printf("[PROTOCOL] data seq=%d epoch=%d path=%s\n",
		pkt.Seq,
		pkt.Epoch,
		pkt.Path,
	)
}

func (e *Engine) StartMigration(candidatePath string) {
	req := e.Protocol.StartRecovery(candidatePath)

	e.Runtime.Trace.Record("protocol_migration_requested", "authority transfer requested", map[string]interface{}{
		"target_path":  req.Path,
		"target_epoch": req.Epoch,
	})

	fmt.Printf("[PROTOCOL] migration requested path=%s target_epoch=%d\n",
		req.Path,
		req.Epoch,
	)
}

func (e *Engine) CommitMigration(candidatePath string) error {
	oldEpoch := e.Protocol.Epoch
	newEpoch := e.Protocol.Epoch + 1

	if err := e.Protocol.ApplyAuthorityTransfer(candidatePath, newEpoch); err != nil {
		return err
	}

	e.Runtime.Trace.Record("protocol_migration_committed", "authority transfer committed", map[string]interface{}{
		"path":  candidatePath,
		"epoch": e.Protocol.Epoch,
	})

	fmt.Printf("[PROTOCOL] migration committed path=%s epoch=%d\n",
		candidatePath,
		e.Protocol.Epoch,
	)

	// protocol-level invariant: epoch must increase
	epochCheck := e.Invariant.CheckEpochMonotonic(oldEpoch, e.Protocol.Epoch)
	e.printInvariant(epochCheck)
	e.recordInvariant(epochCheck)

	return nil
}

func (e *Engine) Receive(pkt protocol.WirePacket) error {
	if err := e.Protocol.ValidatePacket(pkt); err != nil {
		e.Runtime.Trace.Record("protocol_rejected", "packet rejected", map[string]interface{}{
			"reason": err.Error(),
			"seq":    pkt.Seq,
			"epoch":  pkt.Epoch,
			"path":   pkt.Path,
		})
		return err
	}

	e.Runtime.Trace.Record("protocol_accepted", "packet accepted", map[string]interface{}{
		"type":  pkt.Type,
		"seq":   pkt.Seq,
		"epoch": pkt.Epoch,
		"path":  pkt.Path,
	})

	fmt.Printf("[PROTOCOL] accepted type=%s seq=%d epoch=%d path=%s\n",
		pkt.Type,
		pkt.Seq,
		pkt.Epoch,
		pkt.Path,
	)

	// packet-level invariant
	packetCheck := e.Invariant.CheckPacket(pkt, e.Protocol)
	e.printInvariant(packetCheck)
	e.recordInvariant(packetCheck)

	return nil
}

func (e *Engine) CheckProtocolInvariants() {
	fmt.Println("\n=== PROTOCOL INVARIANTS ===")

	results := e.Invariant.RunAll(e.Protocol)
	for _, res := range results {
		e.printInvariant(res)
		e.recordInvariant(res)
	}
}

func (e *Engine) printInvariant(res protocol.InvariantResult) {
	if res.Passed {
		fmt.Printf("[INVARIANT OK] %s — %s\n", res.Name, res.Reason)
		return
	}

	fmt.Printf("[INVARIANT FAIL] %s — %s\n", res.Name, res.Reason)
}

func (e *Engine) recordInvariant(res protocol.InvariantResult) {
	e.Runtime.Trace.Record("protocol_invariant", "protocol invariant evaluated", map[string]interface{}{
		"name":   res.Name,
		"passed": res.Passed,
		"reason": res.Reason,
	})
}