package protocol

import "fmt"

type InvariantResult struct {
	Name   string
	Passed bool
	Reason string
}

type InvariantChecker struct{}

func NewInvariantChecker() *InvariantChecker {
	return &InvariantChecker{}
}

func (ic *InvariantChecker) CheckSessionIdentity(sp *SessionProtocol) InvariantResult {
	if sp == nil {
		return InvariantResult{
			Name:   "session_identity_exists",
			Passed: false,
			Reason: "nil session protocol",
		}
	}

	if sp.SessionID == "" {
		return InvariantResult{
			Name:   "session_identity_exists",
			Passed: false,
			Reason: "empty session id",
		}
	}

	return InvariantResult{
		Name:   "session_identity_exists",
		Passed: true,
		Reason: "session id is present",
	}
}

func (ic *InvariantChecker) CheckEpochMonotonic(oldEpoch, newEpoch int) InvariantResult {
	if newEpoch <= oldEpoch {
		return InvariantResult{
			Name:   "epoch_monotonic",
			Passed: false,
			Reason: fmt.Sprintf("new epoch %d is not greater than old epoch %d", newEpoch, oldEpoch),
		}
	}

	return InvariantResult{
		Name:   "epoch_monotonic",
		Passed: true,
		Reason: fmt.Sprintf("epoch advanced from %d to %d", oldEpoch, newEpoch),
	}
}

func (ic *InvariantChecker) CheckActivePath(sp *SessionProtocol) InvariantResult {
	if sp == nil {
		return InvariantResult{
			Name:   "active_path_exists",
			Passed: false,
			Reason: "nil session protocol",
		}
	}

	if sp.ActivePath == "" {
		return InvariantResult{
			Name:   "active_path_exists",
			Passed: false,
			Reason: "no active path",
		}
	}

	return InvariantResult{
		Name:   "active_path_exists",
		Passed: true,
		Reason: "active path is set",
	}
}

func (ic *InvariantChecker) CheckAuthorityOwner(sp *SessionProtocol) InvariantResult {
	if sp == nil {
		return InvariantResult{
			Name:   "authority_owner_exists",
			Passed: false,
			Reason: "nil session protocol",
		}
	}

	if sp.AuthorityOwner == "" {
		return InvariantResult{
			Name:   "authority_owner_exists",
			Passed: false,
			Reason: "no authority owner",
		}
	}

	if sp.AuthorityOwner != sp.ActivePath {
		return InvariantResult{
			Name:   "authority_matches_active_path",
			Passed: false,
			Reason: fmt.Sprintf("authority owner %s does not match active path %s", sp.AuthorityOwner, sp.ActivePath),
		}
	}

	return InvariantResult{
		Name:   "authority_matches_active_path",
		Passed: true,
		Reason: "authority owner matches active path",
	}
}

func (ic *InvariantChecker) CheckStateAttached(sp *SessionProtocol) InvariantResult {
	if sp == nil {
		return InvariantResult{
			Name:   "session_state_valid",
			Passed: false,
			Reason: "nil session protocol",
		}
	}

	if sp.State != SessionStateAttached && sp.State != SessionStateRecovering && sp.State != SessionStateInit {
		return InvariantResult{
			Name:   "session_state_valid",
			Passed: false,
			Reason: fmt.Sprintf("invalid state %s", sp.State),
		}
	}

	return InvariantResult{
		Name:   "session_state_valid",
		Passed: true,
		Reason: fmt.Sprintf("state is %s", sp.State),
	}
}

func (ic *InvariantChecker) CheckPacket(pkt WirePacket, sp *SessionProtocol) InvariantResult {
	if sp == nil {
		return InvariantResult{
			Name:   "packet_matches_session",
			Passed: false,
			Reason: "nil session protocol",
		}
	}

	if pkt.SessionID != sp.SessionID {
		return InvariantResult{
			Name:   "packet_matches_session",
			Passed: false,
			Reason: "session mismatch",
		}
	}

	if pkt.Epoch != sp.Epoch {
		return InvariantResult{
			Name:   "packet_epoch_current",
			Passed: false,
			Reason: fmt.Sprintf("packet epoch %d != session epoch %d", pkt.Epoch, sp.Epoch),
		}
	}

	return InvariantResult{
		Name:   "packet_epoch_current",
		Passed: true,
		Reason: "packet matches current session epoch",
	}
}

func (ic *InvariantChecker) RunAll(sp *SessionProtocol) []InvariantResult {
	results := []InvariantResult{
		ic.CheckSessionIdentity(sp),
		ic.CheckActivePath(sp),
		ic.CheckAuthorityOwner(sp),
		ic.CheckStateAttached(sp),
	}

	return results
}