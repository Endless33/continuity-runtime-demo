package protocol

import "fmt"

type TransitionResult struct {
	Allowed bool
	Reason  string
	From    SessionState
	To      SessionState
}

type StateMachine struct{}

func NewStateMachine() *StateMachine {
	return &StateMachine{}
}

func (sm *StateMachine) CanTransition(from, to SessionState) TransitionResult {
	switch from {
	case SessionStateInit:
		if to == SessionStateAttached || to == SessionStateClosed {
			return TransitionResult{
				Allowed: true,
				Reason:  "valid init transition",
				From:    from,
				To:      to,
			}
		}

	case SessionStateAttached:
		if to == SessionStateRecovering || to == SessionStateClosed {
			return TransitionResult{
				Allowed: true,
				Reason:  "valid attached transition",
				From:    from,
				To:      to,
			}
		}

	case SessionStateRecovering:
		if to == SessionStateAttached || to == SessionStateClosed {
			return TransitionResult{
				Allowed: true,
				Reason:  "valid recovering transition",
				From:    from,
				To:      to,
			}
		}

	case SessionStateClosed:
		return TransitionResult{
			Allowed: false,
			Reason:  "closed is terminal",
			From:    from,
			To:      to,
		}
	}

	return TransitionResult{
		Allowed: false,
		Reason:  fmt.Sprintf("invalid transition %s -> %s", from, to),
		From:    from,
		To:      to,
	}
}

func (sm *StateMachine) Apply(sp *SessionProtocol, to SessionState) error {
	if sp == nil {
		return fmt.Errorf("nil session protocol")
	}

	res := sm.CanTransition(sp.State, to)
	if !res.Allowed {
		return fmt.Errorf(res.Reason)
	}

	sp.State = to
	return nil
}