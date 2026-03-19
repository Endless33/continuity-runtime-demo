package protocol

import "fmt"

type AuthorityDecision string

const (
	AuthorityAllowTransfer AuthorityDecision = "allow_transfer"
	AuthorityRejectStale   AuthorityDecision = "reject_stale"
	AuthorityRejectPath    AuthorityDecision = "reject_path"
	AuthorityRejectSession AuthorityDecision = "reject_session"
	AuthorityRejectEpoch   AuthorityDecision = "reject_epoch"
)

type AuthorityResult struct {
	Decision AuthorityDecision
	Allowed  bool
	Reason   string
}

type AuthorityRules struct{}

func NewAuthorityRules() *AuthorityRules {
	return &AuthorityRules{}
}

// CanStartTransfer decides whether a migration request is valid.
func (ar *AuthorityRules) CanStartTransfer(sp *SessionProtocol, candidatePath string, targetEpoch int) AuthorityResult {
	if sp == nil {
		return AuthorityResult{
			Decision: AuthorityRejectSession,
			Allowed:  false,
			Reason:   "nil session protocol",
		}
	}

	if candidatePath == "" {
		return AuthorityResult{
			Decision: AuthorityRejectPath,
			Allowed:  false,
			Reason:   "empty candidate path",
		}
	}

	if candidatePath == sp.ActivePath {
		return AuthorityResult{
			Decision: AuthorityRejectPath,
			Allowed:  false,
			Reason:   "candidate path is already active",
		}
	}

	if targetEpoch <= sp.Epoch {
		return AuthorityResult{
			Decision: AuthorityRejectEpoch,
			Allowed:  false,
			Reason:   fmt.Sprintf("target epoch %d must be greater than current epoch %d", targetEpoch, sp.Epoch),
		}
	}

	return AuthorityResult{
		Decision: AuthorityAllowTransfer,
		Allowed:  true,
		Reason:   "transfer allowed",
	}
}

// ValidateIncomingPath checks whether a packet path is still authoritative.
func (ar *AuthorityRules) ValidateIncomingPath(sp *SessionProtocol, pkt WirePacket) AuthorityResult {
	if sp == nil {
		return AuthorityResult{
			Decision: AuthorityRejectSession,
			Allowed:  false,
			Reason:   "nil session protocol",
		}
	}

	if pkt.SessionID != sp.SessionID {
		return AuthorityResult{
			Decision: AuthorityRejectSession,
			Allowed:  false,
			Reason:   "session mismatch",
		}
	}

	if pkt.Epoch < sp.Epoch {
		return AuthorityResult{
			Decision: AuthorityRejectStale,
			Allowed:  false,
			Reason:   fmt.Sprintf("packet epoch %d is stale (current epoch %d)", pkt.Epoch, sp.Epoch),
		}
	}

	if pkt.Epoch > sp.Epoch {
		return AuthorityResult{
			Decision: AuthorityRejectEpoch,
			Allowed:  false,
			Reason:   fmt.Sprintf("packet epoch %d is ahead of current epoch %d", pkt.Epoch, sp.Epoch),
		}
	}

	if pkt.Path != "" && pkt.Path != sp.ActivePath {
		return AuthorityResult{
			Decision: AuthorityRejectPath,
			Allowed:  false,
			Reason:   fmt.Sprintf("packet path %s is not active path %s", pkt.Path, sp.ActivePath),
		}
	}

	return AuthorityResult{
		Decision: AuthorityAllowTransfer,
		Allowed:  true,
		Reason:   "path accepted",
	}
}

// CanCommitTransfer decides whether authority may move to the candidate path.
func (ar *AuthorityRules) CanCommitTransfer(sp *SessionProtocol, candidatePath string, newEpoch int) AuthorityResult {
	if sp == nil {
		return AuthorityResult{
			Decision: AuthorityRejectSession,
			Allowed:  false,
			Reason:   "nil session protocol",
		}
	}

	if candidatePath == "" {
		return AuthorityResult{
			Decision: AuthorityRejectPath,
			Allowed:  false,
			Reason:   "empty candidate path",
		}
	}

	if newEpoch <= sp.Epoch {
		return AuthorityResult{
			Decision: AuthorityRejectEpoch,
			Allowed:  false,
			Reason:   fmt.Sprintf("new epoch %d must be greater than current epoch %d", newEpoch, sp.Epoch),
		}
	}

	return AuthorityResult{
		Decision: AuthorityAllowTransfer,
		Allowed:  true,
		Reason:   "commit allowed",
	}
}