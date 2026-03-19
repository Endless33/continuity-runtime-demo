package main

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
	"continuity-runtime-demo/internal/runtime"
	"continuity-runtime-demo/internal/sim"
)

func main() {
	fmt.Println("TWO NODE PROTOCOL DEMO")
	fmt.Println()

	client := runtime.NewEngine()
	server := runtime.NewEngine()

	client.Init()
	server.Init()

	link := sim.NewLink("client->server")

	fmt.Println("\n=== PHASE 1: SESSION INIT ===")
	initPkt := client.Protocol.BuildInit()
	if err := link.Deliver(initPkt, server.Receive); err != nil {
		fmt.Printf("[ERROR] init delivery failed: %v\n", err)
	}

	fmt.Println("\n=== PHASE 2: DATA PACKETS ===")
	p1 := client.Protocol.BuildData([]byte("payload-1"))
	p2 := client.Protocol.BuildData([]byte("payload-2"))

	_ = link.Deliver(p1, server.Receive)
	_ = link.Deliver(p2, server.Receive)

	fmt.Println("\n=== PHASE 3: FAILURE + MIGRATION ===")
	client.Runtime.HandleEvent(runtime.EventWiFiFailed)

	if err := client.StartMigration(client.Runtime.Current.Name); err != nil {
		fmt.Printf("[ERROR] start migration failed: %v\n", err)
		return
	}

	req := protocol.WirePacket{
		Type:      protocol.PacketTypeAuthorityTransfer,
		SessionID: client.Protocol.SessionID,
		Epoch:     client.Protocol.Epoch + 1,
		Path:      client.Runtime.Current.Name,
	}

	_ = link.Deliver(req, server.Receive)

	if err := client.CommitMigration(client.Runtime.Current.Name); err != nil {
		fmt.Printf("[ERROR] commit migration failed: %v\n", err)
		return
	}

	fmt.Println("\n=== PHASE 4: DATA AFTER MIGRATION ===")
	p3 := client.Protocol.BuildData([]byte("payload-3"))
	_ = link.Deliver(p3, server.Receive)

	fmt.Println("\n=== PHASE 5: CLIENT INVARIANTS ===")
	client.CheckProtocolInvariants()

	fmt.Println("\n=== PHASE 6: SERVER INVARIANTS ===")
	server.CheckProtocolInvariants()
}