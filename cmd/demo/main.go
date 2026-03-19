package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (SESSION + TRACE + MULTIPATH)")
	fmt.Println()

	r := runtime.NewRuntime()
	net := runtime.NewNetworkSimulator()
	stream := runtime.NewStream(r, net)

	fmt.Println("=== BEFORE FAILURE ===")
	stream.Send(5)

	// adaptive / overlap (можешь оставить или убрать — теперь это controlled runtime)
	stream.Multi.StartOverlap()

	fmt.Println("\n=== FAILURE EVENT ===")
	r.HandleEvent(runtime.EventWiFiFailed)

	fmt.Println("\n=== DURING MIGRATION ===")
	stream.Send(5)

	stream.Multi.StopOverlap()

	fmt.Println("\n=== AFTER MIGRATION ===")
	stream.Send(5)

	// 🔥 ВАЖНО: вывод timeline
	r.Trace.PrintTimeline()
}