package runtime

import (
	"encoding/json"
	"fmt"
	"time"
)

type TraceEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"session_id"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type TraceRecorder struct {
	SessionID string
	Events    []TraceEvent
}

func NewTraceRecorder(sessionID string) *TraceRecorder {
	return &TraceRecorder{
		SessionID: sessionID,
		Events:    []TraceEvent{},
	}
}

func (t *TraceRecorder) Record(eventType, message string, data map[string]interface{}) {
	ev := TraceEvent{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		SessionID: t.SessionID,
		Message:   message,
		Data:      data,
	}

	t.Events = append(t.Events, ev)

	// JSONL вывод (очень важно для demo)
	b, _ := json.Marshal(ev)
	fmt.Println(string(b))
}

func (t *TraceRecorder) PrintTimeline() {
	fmt.Println("\n=== TIMELINE ===")

	for i, ev := range t.Events {
		fmt.Printf("[T+%02d] %s → %s\n", i, ev.Type, ev.Message)
	}
}