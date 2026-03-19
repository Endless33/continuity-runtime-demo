package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (ENGINE + PROTOCOL + AUTHORITY)")
	fmt.Println()

	engine := runtime.NewEngine()
	engine.Init()

	r := engine.Runtime
	net := runtime.NewNetworkSimulator()
	stream := runtime.NewStream(r, net)

	// PHASE 1: PROTOCOL DATA
	fmt.Println("\n=== PHASE 1: PROTOCOL DATA BEFORE FAILURE ===")
	engine.SendData([]byte("hello-1"))
	engine.SendData([]byte("hello-2"))
	engine.SendData([]byte("hello-3"))

	// PHASE 2: NORMAL STREAM
	fmt.Println("\n=== PHASE 2: STREAM BEFORE FAILURE ===")
	stream.Send(5)

	// PHASE 3: OVERLAP START
	fmt.Println("\n=== PHASE 3: START OVERLAP ===")
	stream.Multi.StartOverlap()

	// PHASE 4: FAILURE
	fmt.Println("\n=== PHASE 4: FAILURE EVENT ===")
	r.HandleEvent(runtime.EventWiFiFailed)

	// PHASE 5: PROTOCOL MIGRATION
	fmt.Println("\n=== PHASE 5: PROTOCOL MIGRATION ===")

	if err := engine.StartMigration(r.Current.Name); err != nil {
		fmt.Printf("[ERROR] migration start failed: %v\n", err)
		return
	}

	if err := engine.CommitMigration(r.Current.Name); err != nil {
		fmt.Printf("[ERROR] migration commit failed: %v\n", err)
		return
	}

	engine.CheckProtocolInvariants()

	// PHASE 6: STREAM DURING MIGRATION
	fmt.Println("\n=== PHASE 6: STREAM DURING / AFTER MIGRATION ===")
	stream.Send(5)

	// PHASE 7: STOP OVERLAP
	fmt.Println("\n=== PHASE 7: STOP OVERLAP ===")
	stream.Multi.StopOverlap()
	stream.Send(5)

	// PHASE 8: TRACE TIMELINE
	fmt.Println("\n=== PHASE 8: TIMELINE ===")
	r.Trace.PrintTimeline()

	// PHASE 9: REPLAY
	fmt.Println("\n=== PHASE 9: REPLAY ===")
	replay := runtime.NewReplayEngine(r.Trace.Events)
	replay.Run()

	// PHASE 10: RUNTIME INVARIANTS
	fmt.Println("\n=== PHASE 10: RUNTIME INVARIANTS ===")
	checker := runtime.NewInvariantChecker()
	checker.Check(r.Trace.Events)
}