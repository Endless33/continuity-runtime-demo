package protocol

import "fmt"

type ErrorCode string

const (
	ErrUnsupportedVersion ErrorCode = "unsupported_version"
	ErrSessionMismatch    ErrorCode = "session_mismatch"
	ErrStaleEpoch         ErrorCode = "stale_epoch"
	ErrInactivePath       ErrorCode = "inactive_path"
	ErrInvalidSequence    ErrorCode = "invalid_sequence"
	ErrReplayDetected     ErrorCode = "replay_detected"
	ErrAckInvalid         ErrorCode = "invalid_ack"
	ErrProtocolClosed     ErrorCode = "protocol_closed"
)

type ProtocolError struct {
	Code    ErrorCode
	Message string
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewProtocolError(code ErrorCode, msg string) error {
	return ProtocolError{
		Code:    code,
		Message: msg,
	}
}