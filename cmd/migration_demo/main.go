package main

import (
	"fmt"
	"time"

	"continuity-runtime-demo/internal/runtime"
)

type Transport struct {
	Name   string
	Active bool
	Score  float64
}

type Session struct {
	ID         string
	Epoch      uint64
	Transport  *Transport
	Sequence   uint64
	State      string
}

func main() {
	engine := runtime.NewDecisionEngine(runtime.DecisionEngineConfig{
		DegradeRTTThreshold: 120.0,
		FailRTTThreshold:    220.0,
		LossThreshold:       0.15,
		ConsecutiveBadRequired: 3,
		RecoverWindow:       5 * time.Second,
		MigrationMargin:     15.0,
	})

	wifi := &Transport{
		Name:   "wifi",
		Active: true,
		Score:  40.0,
	}

	fiveG := &Transport{
		Name:   "5g",
		Active: true,
		Score:  75.0,
	}

	session := &Session{
		ID:        "sess-7f31a9c2",
		Epoch:     1,
		Transport: wifi,
		Sequence:  0,
		State:     "ATTACHED",
	}

	fmt.Println("CONTINUITY RUNTIME MIGRATION DEMO")
	fmt.Println("---------------------------------")
	fmt.Println()

	fmt.Printf("[SESSION] id=%s epoch=%d state=%s transport=%s\n",
		session.ID,
		session.Epoch,
		session.State,
		session.Transport.Name,
	)
	fmt.Println()

	base := time.Now()

	healthySamples := []runtime.PathSample{
		{RTT: 42, Loss: 0.01, Jitter: 4, HeartbeatMiss: false, Timestamp: base.Add(0 * time.Second)},
		{RTT: 44, Loss: 0.01, Jitter: 5, HeartbeatMiss: false, Timestamp: base.Add(1 * time.Second)},
		{RTT: 46, Loss: 0.02, Jitter: 6, HeartbeatMiss: false, Timestamp: base.Add(2 * time.Second)},
	}

	for _, sample := range healthySamples {
		session.Sequence++
		decision := engine.Observe(sample)

		fmt.Printf("[DATA] tx seq=%d transport=%s payload=hello-%d\n",
			session.Sequence,
			session.Transport.Name,
			session.Sequence,
		)
		fmt.Printf("[SIGNAL] rtt=%.1fms loss=%.2f jitter=%.1f\n",
			sample.RTT,
			sample.Loss,
			sample.Jitter,
		)
		fmt.Printf("[DECISION] state=%s score=%.2f confidence=%.2f migrate=%v\n",
			decision.State,
			decision.Score,
			decision.Confidence,
			decision.Migrate,
		)
		if decision.Reason != "" {
			fmt.Printf("[REASON] %s\n", decision.Reason)
		}
		fmt.Println()

		time.Sleep(350 * time.Millisecond)
	}

	fmt.Println("[EVENT] WiFi begins degrading")
	fmt.Println()

	degradedSamples := []runtime.PathSample{
		{RTT: 145, Loss: 0.08, Jitter: 18, HeartbeatMiss: false, Timestamp: base.Add(3 * time.Second)},
		{RTT: 178, Loss: 0.13, Jitter: 25, HeartbeatMiss: false, Timestamp: base.Add(4 * time.Second)},
		{RTT: 245, Loss: 0.22, Jitter: 41, HeartbeatMiss: true, Timestamp: base.Add(5 * time.Second)},
	}

	var lastDecision runtime.Decision

	for _, sample := range degradedSamples {
		session.Sequence++
		lastDecision = engine.Observe(sample)

		fmt.Printf("[DATA] tx seq=%d transport=%s payload=hello-%d\n",
			session.Sequence,
			session.Transport.Name,
			session.Sequence,
		)
		fmt.Printf("[SIGNAL] rtt=%.1fms loss=%.2f jitter=%.1f heartbeat_miss=%v\n",
			sample.RTT,
			sample.Loss,
			sample.Jitter,
			sample.HeartbeatMiss,
		)
		fmt.Printf("[DECISION] state=%s score=%.2f confidence=%.2f migrate=%v\n",
			lastDecision.State,
			lastDecision.Score,
			lastDecision.Confidence,
			lastDecision.Migrate,
		)
		if lastDecision.Reason != "" {
			fmt.Printf("[REASON] %s\n", lastDecision.Reason)
		}
		fmt.Println()

		time.Sleep(450 * time.Millisecond)
	}

	currentScore := lastDecision.Score
	candidateScore := fiveG.Score
	shouldMigrate := engine.EvaluateCandidate(currentScore, candidateScore, lastDecision.Confidence)

	fmt.Printf("[CANDIDATE] transport=%s score=%.2f current_transport=%s current_score=%.2f\n",
		fiveG.Name,
		candidateScore,
		session.Transport.Name,
		currentScore,
	)
	fmt.Printf("[MIGRATION] candidate_better=%v confidence=%.2f\n",
		shouldMigrate,
		lastDecision.Confidence,
	)
	fmt.Println()

	time.Sleep(500 * time.Millisecond)

	if shouldMigrate {
		oldTransport := session.Transport.Name
		session.Epoch++
		session.Transport = fiveG
		session.State = "RECOVERING"

		fmt.Printf("[AUTHORITY] epoch %d granted to %s\n", session.Epoch, session.Transport.Name)
		fmt.Printf("[STATE] ATTACHED -> %s\n", session.State)
		fmt.Printf("[SWITCH] transport %s -> %s\n", oldTransport, session.Transport.Name)
		fmt.Println()

		time.Sleep(500 * time.Millisecond)

		fmt.Println("[CHECK] stale WiFi rejected")
		session.State = "ATTACHED"
		fmt.Printf("[STATE] RECOVERING -> %s\n", session.State)
		fmt.Println()
	}

	recoveredSamples := []runtime.PathSample{
		{RTT: 68, Loss: 0.02, Jitter: 8, HeartbeatMiss: false, Timestamp: base.Add(6 * time.Second)},
		{RTT: 64, Loss: 0.01, Jitter: 7, HeartbeatMiss: false, Timestamp: base.Add(7 * time.Second)},
		{RTT: 61, Loss: 0.01, Jitter: 6, HeartbeatMiss: false, Timestamp: base.Add(8 * time.Second)},
	}

	for _, sample := range recoveredSamples {
		session.Sequence++
		decision := engine.Observe(sample)

		fmt.Printf("[DATA] tx seq=%d transport=%s payload=hello-%d\n",
			session.Sequence,
			session.Transport.Name,
			session.Sequence,
		)
		fmt.Printf("[SIGNAL] rtt=%.1fms loss=%.2f jitter=%.1f\n",
			sample.RTT,
			sample.Loss,
			sample.Jitter,
		)
		fmt.Printf("[DECISION] state=%s score=%.2f confidence=%.2f migrate=%v\n",
			decision.State,
			decision.Score,
			decision.Confidence,
			decision.Migrate,
		)
		if decision.Reason != "" {
			fmt.Printf("[REASON] %s\n", decision.Reason)
		}
		fmt.Println()

		time.Sleep(350 * time.Millisecond)
	}

	fmt.Println("[RESULT] session continues")
	fmt.Printf("[RESULT] same_session_id=%s\n", session.ID)
	fmt.Printf("[RESULT] current_epoch=%d\n", session.Epoch)
	fmt.Printf("[RESULT] current_transport=%s\n", session.Transport.Name)
	fmt.Println("[RESULT] no reconnect")
	fmt.Println("[RESULT] no session reset")
}