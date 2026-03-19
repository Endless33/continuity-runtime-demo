package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (LOSS + JITTER + STREAM)")
	fmt.Println()

	r := runtime.NewRuntime()
	net := runtime.NewNetworkSimulator()
	stream := runtime.NewStream(r, net)

	// before failure
	stream.Send(5)

	// simulate failure
	r.HandleEvent(runtime.EventWiFiFailed)

	// after migration
	stream.Send(5)
}