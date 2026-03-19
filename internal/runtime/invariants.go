package runtime

import "fmt"

type InvariantChecker struct{}

func NewInvariantChecker() *InvariantChecker {
	return &InvariantChecker{}
}

func (ic *InvariantChecker) Check(events []TraceEvent) {
	fmt.Println("\n=== INVARIANTS ===")

	epoch := 0

	for _, ev := range events {
		if ev.Type == "authority_granted" {
			newEpoch := int(ev.Data["epoch"].(float64))

			if newEpoch <= epoch {
				fmt.Println("[FAIL] epoch not increasing")
				return
			}

			epoch = newEpoch
		}
	}

	fmt.Println("[OK] epoch monotonic")
}