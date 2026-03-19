package main

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("TWO NODE PROTOCOL DEMO (LOSSY EXCHANGE + ACK EXCHANGE + MIGRATION)")
	fmt.Println()

	client := runtime.NewNode("client")
	server := runtime.NewNode("server")

	client.Engine.Init()
	server.Engine.Init()

	forward := runtime.NewLossyExchange("client->server")
	reverse := runtime.NewLossyExchange("server->client")
	acks := runtime.NewAckExchange(forward, reverse)

	fmt.Println("\n=== PHASE 1: SESSION INIT ===")
	initPkt := client.StartHandshake()
	resp, err := forward.Send(client, server, initPkt)
	if err != nil {
		fmt.Printf("[ERROR] init delivery failed: %v\n", err)
		return
	}
	if resp != nil {
		if _, err := reverse.Send(server, client, *resp); err != nil {
			fmt.Printf("[ERROR] init ack delivery failed: %v\n", err)
			return
		}
	}

	fmt.Println("\n=== PHASE 2: DATA PACKETS ===")
	p1 := client.Engine.Protocol.BuildData([]byte("payload-1"))
	p2 := client.Engine.Protocol.BuildData([]byte("payload-2"))

	if err := acks.SendData(client, server, p1); err != nil {
		fmt.Printf("[ERROR] failed to exchange p1: %v\n", err)
	}
	if err := acks.SendData(client, server, p2); err != nil {
		fmt.Printf("[ERROR] failed to exchange p2: %v\n", err)
	}

	fmt.Println("\n=== PHASE 3: FAILURE + MIGRATION ===")
	client.Engine.Runtime.HandleEvent(runtime.EventWiFiFailed)

	if err := client.Engine.StartMigration(client.Engine.Runtime.Current.Name); err != nil {
		fmt.Printf("[ERROR] start migration failed: %v\n", err)
		return
	}

	req := protocol.WirePacket{
		Type:      protocol.PacketTypeAuthorityTransfer,
		SessionID: client.Engine.Protocol.SessionID,
		Epoch:     client.Engine.Protocol.Epoch + 1,
		Path:      client.Engine.Runtime.Current.Name,
	}

	if _, err := forward.Send(client, server, req); err != nil {
		fmt.Printf("[ERROR] transfer request delivery failed: %v\n", err)
	}

	if err := client.Engine.CommitMigration(client.Engine.Runtime.Current.Name); err != nil {
		fmt.Printf("[ERROR] commit migration failed: %v\n", err)
		return
	}

	fmt.Println("\n=== PHASE 4: DATA AFTER MIGRATION ===")
	p3 := client.Engine.Protocol.BuildData([]byte("payload-3"))
	if err := acks.SendData(client, server, p3); err != nil {
		fmt.Printf("[ERROR] failed to exchange p3: %v\n", err)
	}

	fmt.Println("\n=== PHASE 5: CLIENT INVARIANTS ===")
	client.Engine.CheckProtocolInvariants()

	fmt.Println("\n=== PHASE 6: SERVER INVARIANTS ===")
	server.Engine.CheckProtocolInvariants()
}