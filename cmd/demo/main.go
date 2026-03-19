package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (STREAM + LOSS + FEC + MIGRATION)")
	fmt.Println()

	r := runtime.NewRuntime()
	net := runtime.NewNetworkSimulator()
	stream := runtime.NewStream(r, net)

	fmt.Println("=== PHASE 1: NORMAL STREAMING ===")
	stream.Send(5)

	fmt.Println("\n=== PHASE 2: OVERLAP ENABLED ===")
	stream.Multi.StartOverlap()
	stream.Send(3)

	fmt.Println("\n=== PHASE 3: FAILURE EVENT ===")
	r.HandleEvent(runtime.EventWiFiFailed)

	fmt.Println("\n=== PHASE 4: STREAM DURING MIGRATION ===")
	stream.Send(5)

	fmt.Println("\n=== PHASE 5: OVERLAP DISABLED ===")
	stream.Multi.StopOverlap()
	stream.Send(5)

	fmt.Println("\n=== PHASE 6: TIMELINE ===")
	r.Trace.PrintTimeline()

	fmt.Println("\n=== PHASE 7: REPLAY ===")
	replay := runtime.NewReplayEngine(r.Trace.Events)
	replay.Run()

	fmt.Println("\n=== PHASE 8: INVARIANTS ===")
	checker := runtime.NewInvariantChecker()
	checker.Check(r.Trace.Events)
}