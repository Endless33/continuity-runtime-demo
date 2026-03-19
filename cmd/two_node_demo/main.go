package main

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("TWO NODE PROTOCOL DEMO (LOSSY EXCHANGE + ACK + MIGRATION)")
	fmt.Println()

	client := runtime.NewNode("client")
	server := runtime.NewNode("server")

	client.Engine.Init()
	server.Engine.Init()

	forward := runtime.NewLossyExchange("client->server")
	reverse := runtime.NewLossyExchange("server->client")
	ackRules := protocol.NewAckRules()

	serverReceive := func(pkt protocol.WirePacket) error {
		_, err := server.Receive(pkt)
		if err != nil {
			fmt.Printf("[SERVER ERROR] %v\n", err)
			return err
		}

		if pkt.Type == protocol.PacketTypeData {
			ack := server.Engine.Protocol.BuildAck(pkt.Seq)

			if _, err := reverse.Send(server, client, ack); err != nil {
				fmt.Printf("[CLIENT ERROR] %v\n", err)
				return err
			}

			ackCheck := ackRules.ValidateAck(client.Engine.Protocol, ack)
			if ackCheck.Allowed {
				fmt.Printf("[ACK RULES] accepted ack=%d reason=%s\n", ack.Ack, ackCheck.Reason)
			} else {
				fmt.Printf("[ACK RULES] rejected ack=%d reason=%s\n", ack.Ack, ackCheck.Reason)
			}
		}

		return nil
	}

	fmt.Println("\n=== PHASE 1: SESSION INIT ===")
	initPkt := client.StartHandshake()
	if _, err := forward.Send(client, server, initPkt); err != nil {
		fmt.Printf("[ERROR] init delivery failed: %v\n", err)
	}

	fmt.Println("\n=== PHASE 2: DATA PACKETS ===")
	p1 := client.Engine.Protocol.BuildData([]byte("payload-1"))
	p2 := client.Engine.Protocol.BuildData([]byte("payload-2"))

	_ = serverReceive(p1)
	_ = serverReceive(p2)

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
	_ = serverReceive(p3)

	fmt.Println("\n=== PHASE 5: CLIENT INVARIANTS ===")
	client.Engine.CheckProtocolInvariants()

	fmt.Println("\n=== PHASE 6: SERVER INVARIANTS ===")
	server.Engine.CheckProtocolInvariants()
}