package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (STRUCTURED)")
	fmt.Println()

	r := runtime.NewRuntime()

	// before failure
	for i := 0; i < 3; i++ {
		r.SendPacket()
	}

	// simulate failure
	r.HandleEvent(runtime.EventWiFiFailed)

	// after migration
	for i := 0; i < 3; i++ {
		r.SendPacket()
	}
}