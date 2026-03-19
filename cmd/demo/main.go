package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (SESSION + TRACE + MULTIPATH + REPLAY)")
	fmt.Println()

	r := runtime.NewRuntime()
	net := runtime.NewNetworkSimulator()
	stream := runtime.NewStream(r, net)

	fmt.Println("=== BEFORE FAILURE ===")
	stream.Send(5)

	stream.Multi.StartOverlap()

	fmt.Println("\n=== FAILURE EVENT ===")
	r.HandleEvent(runtime.EventWiFiFailed)

	fmt.Println("\n=== DURING MIGRATION ===")
	stream.Send(5)

	stream.Multi.StopOverlap()

	fmt.Println("\n=== AFTER MIGRATION ===")
	stream.Send(5)

	r.Trace.PrintTimeline()

	replay := runtime.NewReplayEngine(r.Trace.Events)
	replay.Run()

	checker := runtime.NewInvariantChecker()
	checker.Check(r.Trace.Events)
}