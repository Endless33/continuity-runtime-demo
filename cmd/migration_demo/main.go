package main

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("MIGRATION DEMO")
	fmt.Println()

	client := runtime.NewNode("client")
	server := runtime.NewNode("server")

	// Handshake first
	initPkt := client.StartHandshake()
	ackPkt, err := server.Receive(initPkt)
	if err != nil {
		fmt.Printf("[ERROR] server failed to handle init: %v\n", err)
		return
	}
	if ackPkt == nil {
		fmt.Println("[ERROR] missing init ack")
		return
	}
	if _, err := client.Receive(*ackPkt); err != nil {
		fmt.Printf("[ERROR] client failed to handle init ack: %v\n", err)
		return
	}

	fmt.Println("\n=== DATA BEFORE FAILURE ===")
	client.Engine.SendData([]byte("pre-failure-1"))
	client.Engine.SendData([]byte("pre-failure-2"))

	fmt.Println("\n=== FAILURE + RECOVERY START ===")
	client.Engine.Runtime.HandleEvent(runtime.EventWiFiFailed)

	if err := client.BeginRecovery(); err != nil {
		fmt.Printf("[ERROR] client could not enter recovering: %v\n", err)
		return
	}
	if err := server.BeginRecovery(); err != nil {
		fmt.Printf("[ERROR] server could not enter recovering: %v\n", err)
		return
	}

	targetPath := client.Engine.Runtime.Current.Name

	fmt.Println("\n=== MIGRATION REQUEST ===")
	if err := client.Engine.StartMigration(targetPath); err != nil {
		fmt.Printf("[ERROR] client start migration failed: %v\n", err)
		return
	}

	req := protocol.WirePacket{
		Type:      protocol.PacketTypeAuthorityTransfer,
		SessionID: client.Engine.Protocol.SessionID,
		Epoch:     client.Engine.Protocol.Epoch + 1,
		Path:      targetPath,
	}

	if _, err := server.Receive(req); err != nil {
		fmt.Printf("[ERROR] server rejected transfer request: %v\n", err)
	}

	fmt.Println("\n=== SERVER COMMITS MIGRATION ===")
	if err := server.Engine.Protocol.ApplyAuthorityTransfer(targetPath, req.Epoch); err != nil {
		fmt.Printf("[ERROR] server apply transfer failed: %v\n", err)
		return
	}
	if err := server.EndRecovery(); err != nil {
		fmt.Printf("[ERROR] server could not return attached: %v\n", err)
		return
	}

	fmt.Printf("[SERVER] authority path=%s epoch=%d\n",
		server.Engine.Protocol.ActivePath,
		server.Engine.Protocol.Epoch,
	)

	fmt.Println("\n=== CLIENT COMMITS MIGRATION ===")
	if err := client.Engine.CommitMigration(targetPath); err != nil {
		fmt.Printf("[ERROR] client commit migration failed: %v\n", err)
		return
	}
	if err := client.EndRecovery(); err != nil {
		fmt.Printf("[ERROR] client could not return attached: %v\n", err)
		return
	}

	client.Engine.CheckProtocolInvariants()
	server.Engine.CheckProtocolInvariants()

	fmt.Println("\n=== DATA AFTER MIGRATION ===")
	client.Engine.SendData([]byte("post-migration-1"))
	client.Engine.SendData([]byte("post-migration-2"))

	fmt.Println("\n=== RESULT ===")
	fmt.Printf("[CLIENT] state=%s epoch=%d path=%s\n",
		client.Engine.Protocol.State,
		client.Engine.Protocol.Epoch,
		client.Engine.Protocol.ActivePath,
	)

	fmt.Printf("[SERVER] state=%s epoch=%d path=%s\n",
		server.Engine.Protocol.State,
		server.Engine.Protocol.Epoch,
		server.Engine.Protocol.ActivePath,
	)

	fmt.Println("\n=== CLIENT TIMELINE ===")
	client.Engine.Runtime.Trace.PrintTimeline()

	fmt.Println("\n=== SERVER TIMELINE ===")
	server.Engine.Runtime.Trace.PrintTimeline()
}