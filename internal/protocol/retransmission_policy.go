package protocol

import "time"

type RetransmissionPolicy struct {
	MaxRetries      int
	BaseBackoff     time.Duration
	MaxBackoff      time.Duration
	DuplicateAckGap int
}

func NewRetransmissionPolicy() *RetransmissionPolicy {
	return &RetransmissionPolicy{
		MaxRetries:      5,
		BaseBackoff:     200 * time.Millisecond,
		MaxBackoff:      2 * time.Second,
		DuplicateAckGap: 3,
	}
}

func (rp *RetransmissionPolicy) BackoffFor(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}

	backoff := rp.BaseBackoff
	for i := 1; i < attempt; i++ {
		backoff *= 2
		if backoff >= rp.MaxBackoff {
			return rp.MaxBackoff
		}
	}

	if backoff > rp.MaxBackoff {
		return rp.MaxBackoff
	}

	return backoff
}