package runtime

import (
	"fmt"
	"math"
	"time"
)

type PathState int

const (
	HEALTHY PathState = iota
	DEGRADED
	FAILED
)

func (s PathState) String() string {
	switch s {
	case HEALTHY:
		return "HEALTHY"
	case DEGRADED:
		return "DEGRADED"
	case FAILED:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

type PathSample struct {
	RTT            float64   // ms
	Loss           float64   // 0.0 - 1.0
	Jitter         float64   // ms
	HeartbeatMiss  bool
	Timestamp      time.Time
}

type Decision struct {
	State      PathState
	Migrate    bool
	Score      float64
	Confidence float64
	Reason     string
}

type EWMA struct {
	alpha float64
	value float64
	init  bool
}

func NewEWMA(alpha float64) *EWMA {
	return &EWMA{alpha: alpha}
}

func (e *EWMA) Update(v float64) float64 {
	if !e.init {
		e.value = v
		e.init = true
		return e.value
	}
	e.value = e.alpha*v + (1.0-e.alpha)*e.value
	return e.value
}

type DecisionEngineConfig struct {
	DegradeRTTThreshold float64
	FailRTTThreshold    float64

	LossThreshold float64

	ConsecutiveBadRequired int

	RecoverWindow time.Duration

	MigrationMargin float64
}

type PathMetrics struct {
	rttFast  *EWMA
	rttSlow  *EWMA
	loss     *EWMA
	jitter   *EWMA

	consecutiveBad int
	lastBadTime    time.Time
	lastGoodTime   time.Time

	state PathState
}

type DecisionEngine struct {
	cfg     DecisionEngineConfig
	metrics *PathMetrics
}

func NewDecisionEngine(cfg DecisionEngineConfig) *DecisionEngine {
	return &DecisionEngine{
		cfg: cfg,
		metrics: &PathMetrics{
			rttFast: NewEWMA(0.5),
			rttSlow: NewEWMA(0.1),
			loss:    NewEWMA(0.2),
			jitter:  NewEWMA(0.2),
			state:   HEALTHY,
		},
	}
}

func (e *DecisionEngine) Observe(s PathSample) Decision {
	m := e.metrics

	rttFast := m.rttFast.Update(s.RTT)
	rttSlow := m.rttSlow.Update(s.RTT)
	loss := m.loss.Update(s.Loss)
	_ = m.jitter.Update(s.Jitter)

	now := s.Timestamp

	isBad := false

	if s.HeartbeatMiss {
		isBad = true
	} else if rttFast > e.cfg.DegradeRTTThreshold {
		isBad = true
	} else if loss > e.cfg.LossThreshold {
		isBad = true
	}

	if isBad {
		m.consecutiveBad++
		m.lastBadTime = now
	} else {
		m.consecutiveBad = 0
		m.lastGoodTime = now
	}

	prevState := m.state

	// STATE TRANSITIONS

	if m.consecutiveBad >= e.cfg.ConsecutiveBadRequired {
		if rttFast > e.cfg.FailRTTThreshold || s.HeartbeatMiss {
			m.state = FAILED
		} else {
			m.state = DEGRADED
		}
	} else {
		// recovery via time window
		if now.Sub(m.lastBadTime) > e.cfg.RecoverWindow {
			m.state = HEALTHY
		}
	}

	score := computeScore(rttSlow, loss)

	confidence := computeConfidence(m.consecutiveBad)

	migrate := false
	reason := ""

	if m.state == DEGRADED || m.state == FAILED {
		if confidence > 0.7 {
			migrate = true
			reason = "path degraded with high confidence"
		}
	}

	if prevState != m.state {
		reason = fmt.Sprintf("state transition %s → %s", prevState, m.state)
	}

	return Decision{
		State:      m.state,
		Migrate:    migrate,
		Score:      score,
		Confidence: confidence,
		Reason:     reason,
	}
}

func (e *DecisionEngine) EvaluateCandidate(currentScore, candidateScore float64, confidence float64) bool {
	if candidateScore < currentScore-e.cfg.MigrationMargin && confidence > 0.7 {
		return true
	}
	return false
}

func computeScore(rtt, loss float64) float64 {
	// lower is better
	return rtt + loss*1000.0
}

func computeConfidence(consecutiveBad int) float64 {
	// sigmoid-like growth
	x := float64(consecutiveBad)
	return 1.0 - math.Exp(-x/3.0)
}