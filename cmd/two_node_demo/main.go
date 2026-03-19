package main

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
	"continuity-runtime-demo/internal/runtime"
	"continuity-runtime-demo/internal/sim"
)

func main() {
	fmt.Println("TWO NODE PROTOCOL DEMO (LOSSY LINK + ACK + MIGRATION)")
	fmt.Println()

	client := runtime.NewEngine()
	server := runtime.NewEngine()

	client.Init()
	server.Init()

	forward := sim.NewLossyLink("client->server")
	reverse := sim.NewLossyLink("server->client")
	ackRules := protocol.NewAckRules()

	serverReceive := func(pkt protocol.WirePacket) error {
		if err := server.Receive(pkt); err != nil {
			fmt.Printf("[SERVER ERROR] %v\n", err)
			return err
		}

		// Send ACKs back for data packets
		if pkt.Type == protocol.PacketTypeData {
			ack := server.Protocol.BuildAck(pkt.Seq)

			if err := reverse.Deliver(ack, func(ap protocol.WirePacket) error {
				if err := client.Receive(ap); err != nil {
					fmt.Printf("[CLIENT ERROR] %v\n", err)
					return err
				}

				ackCheck := ackRules.ValidateAck(client.Protocol, ap)
				if ackCheck.Allowed {
					fmt.Printf("[ACK RULES] accepted ack=%d reason=%s\n", ap.Ack, ackCheck.Reason)
				} else {
					fmt.Printf("[ACK RULES] rejected ack=%d reason=%s\n", ap.Ack, ackCheck.Reason)
				}

				return nil
			}); err != nil {
				return err
			}
		}

		return nil
	}

	fmt.Println("\n=== PHASE 1: SESSION INIT ===")
	initPkt := client.Protocol.BuildInit()
	if err := forward.Deliver(initPkt, serverReceive); err != nil {
		fmt.Printf("[ERROR] init delivery failed: %v\n", err)
	}

	fmt.Println("\n=== PHASE 2: DATA PACKETS BEFORE FAILURE ===")
	p1 := client.Protocol.BuildData([]byte("payload-1"))
	p2 := client.Protocol.BuildData([]byte("payload-2"))
	p3 := client.Protocol.BuildData([]byte("payload-3"))

	_ = forward.Deliver(p1, serverReceive)
	_ = forward.Deliver(p2, serverReceive)
	_ = forward.Deliver(p3, serverReceive)

	fmt.Println("\n=== PHASE 3: FAILURE EVENT ON CLIENT ===")
	client.Runtime.HandleEvent(runtime.EventWiFiFailed)

	fmt.Println("\n=== PHASE 4: START MIGRATION ===")
	targetPath := client.Runtime.Current.Name

	if err := client.StartMigration(targetPath); err != nil {
		fmt.Printf("[ERROR] client start migration failed: %v\n", err)
		return
	}

	req := protocol.WirePacket{
		Type:      protocol.PacketTypeAuthorityTransfer,
		SessionID: client.Protocol.SessionID,
		Epoch:     client.Protocol.Epoch + 1,
		Path:      targetPath,
	}

	if err := forward.Deliver(req, serverReceive); err != nil {
		fmt.Printf("[ERROR] transfer request delivery failed: %v\n", err)
		return
	}

	fmt.Println("\n=== PHASE 5: SERVER APPLIES TRANSFER ===")
	if err := server.Protocol.ApplyAuthorityTransfer(targetPath, req.Epoch); err != nil {
		fmt.Printf("[ERROR] server apply transfer failed: %v\n", err)
		return
	}
	fmt.Printf("[SERVER] authority moved path=%s epoch=%d\n", server.Protocol.ActivePath, server.Protocol.Epoch)

	fmt.Println("\n=== PHASE 6: CLIENT COMMITS TRANSFER ===")
	if err := client.CommitMigration(targetPath); err != nil {
		fmt.Printf("[ERROR] client commit migration failed: %v\n", err)
		return
	}

	fmt.Println("\n=== PHASE 7: DATA PACKETS AFTER MIGRATION ===")
	p4 := client.Protocol.BuildData([]byte("payload-4"))
	p5 := client.Protocol.BuildData([]byte("payload-5"))

	_ = forward.Deliver(p4, serverReceive)
	_ = forward.Deliver(p5, serverReceive)

	fmt.Println("\n=== PHASE 8: CLIENT PROTOCOL INVARIANTS ===")
	client.CheckProtocolInvariants()

	fmt.Println("\n=== PHASE 9: SERVER PROTOCOL INVARIANTS ===")
	server.CheckProtocolInvariants()

	fmt.Println("\n=== PHASE 10: CLIENT TRACE TIMELINE ===")
	client.Runtime.Trace.PrintTimeline()

	fmt.Println("\n=== PHASE 11: SERVER TRACE TIMELINE ===")
	server.Runtime.Trace.PrintTimeline()
}