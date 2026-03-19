package runtime

import (
	"fmt"
	"time"
)

func RaceTransports(network *NetworkSimulator, packetID int, a, b Transport) (Transport, bool) {
	type result struct {
		t       Transport
		success bool
		time    time.Duration
	}

	ch := make(chan result, 2)

	send := func(t Transport) {
		start := time.Now()
		ok := network.Transmit(packetID, t)
		ch <- result{
			t:       t,
			success: ok,
			time:    time.Since(start),
		}
	}

	go send(a)
	go send(b)

	r1 := <-ch
	r2 := <-ch

	// выбираем быстрее
	if r1.success && (!r2.success || r1.time <= r2.time) {
		fmt.Printf("[RACE] winner=%s (%v)\n", r1.t.Name, r1.time)
		return r1.t, true
	}

	if r2.success {
		fmt.Printf("[RACE] winner=%s (%v)\n", r2.t.Name, r2.time)
		return r2.t, true
	}

	return a, false
}