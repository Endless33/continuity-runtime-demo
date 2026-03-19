package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO (MULTIPATH + ZERO-LOSS)")
	fmt.Println()

	r := runtime.NewRuntime()
	net := runtime.NewNetworkSimulator()
	stream := runtime.NewStream(r, net)

	// before failure
	stream.Send(5)

	// начинаем overlap ДО падения (важно!)
	stream.Multi.StartOverlap()

	// simulate failure
	r.HandleEvent(runtime.EventWiFiFailed)

	// во время overlap продолжаем стрим
	stream.Send(5)

	// отключаем overlap
	stream.Multi.StopOverlap()

	// продолжаем уже на новом пути
	stream.Send(5)
}