package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (ENGINE + PROTOCOL + STREAM)")
	fmt.Println()

	engine := runtime.NewEngine()
	engine.Init()

	r := engine.Runtime
	net := runtime.NewNetworkSimulator()
	stream := runtime.NewStream(r, net)

	fmt.Println("\n=== PHASE 1: PROTOCOL DATA BEFORE FAILURE ===")
	engine.SendData([]byte("hello-1"))
	engine.SendData([]byte("hello-2"))
	engine.SendData([]byte("hello-3"))

	fmt.Println("\n=== PHASE 2: STREAM BEFORE FAILURE ===")
	stream.Send(5)

	fmt.Println("\n=== PHASE 3: START OVERLAP ===")
	stream.Multi.StartOverlap()

	fmt.Println("\n=== PHASE 4: FAILURE EVENT ===")
	r.HandleEvent(runtime.EventWiFiFailed)

	fmt.Println("\n=== PHASE 5: PROTOCOL MIGRATION ===")
	engine.StartMigration(r.Current.Name)
	if err := engine.CommitMigration(r.Current.Name); err != nil {
		fmt.Printf("[ERROR] migration commit failed: %v\n", err)
		return
	}

	fmt.Println("\n=== PHASE 6: STREAM DURING / AFTER MIGRATION ===")
	stream.Send(5)

	fmt.Println("\n=== PHASE 7: STOP OVERLAP ===")
	stream.Multi.StopOverlap()
	stream.Send(5)

	fmt.Println("\n=== PHASE 8: TIMELINE ===")
	r.Trace.PrintTimeline()

	fmt.Println("\n=== PHASE 9: REPLAY ===")
	replay := runtime.NewReplayEngine(r.Trace.Events)
	replay.Run()

	fmt.Println("\n=== PHASE 10: INVARIANTS ===")
	checker := runtime.NewInvariantChecker()
	checker.Check(r.Trace.Events)
}