package runtime

import (
	"fmt"
	"time"

	"continuity-runtime-demo/internal/protocol"
)

type Engine struct {
	Runtime    *Runtime
	Protocol   *protocol.SessionProtocol
	Invariant  *protocol.InvariantChecker
	Authority  *protocol.AuthorityRules
	Timeouts   *protocol.TimeoutPolicy
	Retransmit *protocol.RetransmissionPolicy
	Keepalive  *protocol.KeepaliveEngine
	Close      *protocol.CloseEngine
}

func NewEngine() *Engine {
	rt := NewRuntime()
	timeouts := protocol.NewTimeoutPolicy()

	return &Engine{
		Runtime:    rt,
		Protocol:   protocol.NewSessionProtocol("sess-001", rt.Current.Name),
		Invariant:  protocol.NewInvariantChecker(),
		Authority:  protocol.NewAuthorityRules(),
		Timeouts:   timeouts,
		Retransmit: protocol.NewRetransmissionPolicy(),
		Keepalive:  protocol.NewKeepaliveEngine(timeouts.KeepaliveInterval),
		Close:      protocol.NewCloseEngine(),
	}
}

func (e *Engine) Init() {
	pkt := e.Protocol.BuildInit()

	e.Runtime.Trace.Record("protocol_init", "session protocol initialized", map[string]interface{}{
		"session_id": pkt.SessionID,
		"epoch":      pkt.Epoch,
		"path":       pkt.Path,
		"version":    pkt.Version,
		"timeout_ms": e.Timeouts.InitTimeout.Milliseconds(),
	})

	fmt.Printf("[PROTOCOL] init session=%s epoch=%d path=%s version=%d\n",
		pkt.SessionID,
		pkt.Epoch,
		pkt.Path,
		pkt.Version,
	)
}

func (e *Engine) SendData(payload []byte) protocol.WirePacket {
	pkt := e.Protocol.BuildData(payload)

	e.Runtime.Trace.Record("protocol_data", "data packet built", map[string]interface{}{
		"seq":            pkt.Seq,
		"epoch":          pkt.Epoch,
		"path":           pkt.Path,
		"version":        pkt.Version,
		"payload_size":   len(pkt.Payload),
		"ack_timeout_ms": e.Timeouts.AckTimeout.Milliseconds(),
		"max_retries":    e.Retransmit.MaxRetries,
	})

	fmt.Printf("[PROTOCOL] data seq=%d epoch=%d path=%s version=%d\n",
		pkt.Seq,
		pkt.Epoch,
		pkt.Path,
		pkt.Version,
	)

	return pkt
}

func (e *Engine) StartMigration(candidatePath string) error {
	targetEpoch := e.Protocol.Epoch + 1

	auth := e.Authority.CanStartTransfer(e.Protocol, candidatePath, targetEpoch)

	e.Runtime.Trace.Record("protocol_authority_start_check", "start transfer evaluated", map[string]interface{}{
		"candidate_path":       candidatePath,
		"target_epoch":         targetEpoch,
		"allowed":              auth.Allowed,
		"decision":             auth.Decision,
		"reason":               auth.Reason,
		"migration_timeout_ms": e.Timeouts.MigrationTimeout.Milliseconds(),
	})

	if !auth.Allowed {
		fmt.Printf("[PROTOCOL] migration request rejected decision=%s reason=%s\n",
			auth.Decision,
			auth.Reason,
		)
		return fmt.Errorf(auth.Reason)
	}

	req := e.Protocol.StartRecovery(candidatePath)

	e.Runtime.Trace.Record("protocol_migration_requested", "authority transfer requested", map[string]interface{}{
		"target_path":  req.Path,
		"target_epoch": req.Epoch,
		"version":      req.Version,
	})

	fmt.Printf("[PROTOCOL] migration requested path=%s target_epoch=%d\n",
		req.Path,
		req.Epoch,
	)

	return nil
}

func (e *Engine) CommitMigration(candidatePath string) error {
	oldEpoch := e.Protocol.Epoch
	newEpoch := e.Protocol.Epoch + 1

	auth := e.Authority.CanCommitTransfer(e.Protocol, candidatePath, newEpoch)

	e.Runtime.Trace.Record("protocol_authority_commit_check", "commit transfer evaluated", map[string]interface{}{
		"path":     candidatePath,
		"epoch":    newEpoch,
		"allowed":  auth.Allowed,
		"decision": auth.Decision,
		"reason":   auth.Reason,
	})

	if !auth.Allowed {
		fmt.Printf("[PROTOCOL] migration commit rejected decision=%s reason=%s\n",
			auth.Decision,
			auth.Reason,
		)
		return fmt.Errorf(auth.Reason)
	}

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

	epochCheck := e.Invariant.CheckEpochMonotonic(oldEpoch, e.Protocol.Epoch)
	e.printInvariant(epochCheck)
	e.recordInvariant(epochCheck)

	return nil
}

func (e *Engine) Receive(pkt protocol.WirePacket) error {
	auth := e.Authority.ValidateIncomingPath(e.Protocol, pkt)

	e.Runtime.Trace.Record("protocol_authority_receive_check", "incoming path evaluated", map[string]interface{}{
		"type":     pkt.Type,
		"seq":      pkt.Seq,
		"ack":      pkt.Ack,
		"epoch":    pkt.Epoch,
		"path":     pkt.Path,
		"version":  pkt.Version,
		"allowed":  auth.Allowed,
		"decision": auth.Decision,
		"reason":   auth.Reason,
	})

	if !auth.Allowed {
		e.Runtime.Trace.Record("protocol_rejected", "packet rejected by authority rules", map[string]interface{}{
			"reason":   auth.Reason,
			"decision": auth.Decision,
			"type":     pkt.Type,
			"seq":      pkt.Seq,
			"ack":      pkt.Ack,
			"epoch":    pkt.Epoch,
			"path":     pkt.Path,
			"version":  pkt.Version,
		})
		return fmt.Errorf(auth.Reason)
	}

	if err := e.Protocol.ValidatePacket(pkt); err != nil {
		e.Runtime.Trace.Record("protocol_rejected", "packet rejected by protocol validation", map[string]interface{}{
			"reason":  err.Error(),
			"type":    pkt.Type,
			"seq":     pkt.Seq,
			"ack":     pkt.Ack,
			"epoch":   pkt.Epoch,
			"path":    pkt.Path,
			"version": pkt.Version,
		})
		return err
	}

	if pkt.Type == protocol.PacketTypeKeepalive {
		e.Keepalive.MarkSeen(time.Now().UTC())
	}

	if pkt.Type == protocol.PacketTypeClose {
		e.Close.Apply(e.Protocol)
	}

	e.Runtime.Trace.Record("protocol_accepted", "packet accepted", map[string]interface{}{
		"type":    pkt.Type,
		"seq":     pkt.Seq,
		"ack":     pkt.Ack,
		"epoch":   pkt.Epoch,
		"path":    pkt.Path,
		"version": pkt.Version,
	})

	fmt.Printf("[PROTOCOL] accepted type=%s seq=%d ack=%d epoch=%d path=%s version=%d\n",
		pkt.Type,
		pkt.Seq,
		pkt.Ack,
		pkt.Epoch,
		pkt.Path,
		pkt.Version,
	)

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

func (e *Engine) SimulateAckWait(seq int) {
	timeout := e.Timeouts.AckTimeout

	e.Runtime.Trace.Record("protocol_ack_wait", "waiting for ack", map[string]interface{}{
		"seq":        seq,
		"timeout_ms": timeout.Milliseconds(),
	})

	fmt.Printf("[PROTOCOL] waiting ack for seq=%d timeout=%v\n", seq, timeout)
	time.Sleep(timeout / 10)
}

func (e *Engine) SimulateRetransmission(seq int, attempt int) {
	backoff := e.Retransmit.BackoffFor(attempt)

	e.Runtime.Trace.Record("protocol_retransmit", "retransmission scheduled", map[string]interface{}{
		"seq":        seq,
		"attempt":    attempt,
		"backoff_ms": backoff.Milliseconds(),
	})

	fmt.Printf("[PROTOCOL] retransmit seq=%d attempt=%d backoff=%v\n",
		seq,
		attempt,
		backoff,
	)
	time.Sleep(backoff / 10)
}

func (e *Engine) SendKeepalive() protocol.WirePacket {
	now := time.Now().UTC()

	if !e.Keepalive.ShouldSend(now) {
		return protocol.WirePacket{}
	}

	pkt := e.Keepalive.Build(e.Protocol)
	e.Keepalive.MarkSent(now)

	e.Runtime.Trace.Record("protocol_keepalive_sent", "keepalive sent", map[string]interface{}{
		"epoch":   pkt.Epoch,
		"path":    pkt.Path,
		"version": pkt.Version,
	})

	fmt.Printf("[PROTOCOL] keepalive sent epoch=%d path=%s\n", pkt.Epoch, pkt.Path)
	return pkt
}

func (e *Engine) CheckIdleTimeout() bool {
	expired := e.Keepalive.IsExpired(time.Now().UTC(), e.Timeouts.IdleSessionTimeout)

	e.Runtime.Trace.Record("protocol_idle_check", "idle timeout evaluated", map[string]interface{}{
		"expired":         expired,
		"idle_timeout_ms": e.Timeouts.IdleSessionTimeout.Milliseconds(),
	})

	if expired {
		fmt.Println("[PROTOCOL] idle timeout expired")
	}

	return expired
}

func (e *Engine) CloseSession(reason protocol.CloseReason) protocol.WirePacket {
	pkt := e.Close.Build(e.Protocol, reason)
	e.Close.Apply(e.Protocol)

	e.Runtime.Trace.Record("protocol_close", "session closed", map[string]interface{}{
		"reason":  reason,
		"epoch":   pkt.Epoch,
		"path":    pkt.Path,
		"version": pkt.Version,
	})

	fmt.Printf("[PROTOCOL] close session reason=%s epoch=%d path=%s\n",
		reason,
		pkt.Epoch,
		pkt.Path,
	)

	return pkt
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