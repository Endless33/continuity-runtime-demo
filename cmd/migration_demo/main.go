package main

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("MIGRATION DEMO (LOSSY EXCHANGE + ACK EXCHANGE)")
	fmt.Println()

	client := runtime.NewNode("client")
	server := runtime.NewNode("server")

	forward := runtime.NewLossyExchange("client->server")
	reverse := runtime.NewLossyExchange("server->client")
	acks := runtime.NewAckExchange(forward, reverse)

	// Handshake
	initPkt := client.StartHandshake()
	resp, err := forward.Send(client, server, initPkt)
	if err != nil {
		fmt.Printf("[ERROR] server failed to handle init: %v\n", err)
		return
	}
	if resp == nil {
		fmt.Println("[WARN] init dropped or missing init ack")
		return
	}
	if _, err := reverse.Send(server, client, *resp); err != nil {
		fmt.Printf("[ERROR] client failed to handle init ack: %v\n", err)
		return
	}

	fmt.Println("\n=== DATA BEFORE FAILURE ===")
	client.Engine.SendData([]byte("pre-failure-1"))
	client.Engine.SendData([]byte("pre-failure-2"))

	data1 := client.Engine.Protocol.BuildData([]byte("wire-pre-1"))
	data2 := client.Engine.Protocol.BuildData([]byte("wire-pre-2"))

	if err := acks.SendData(client, server, data1); err != nil {
		fmt.Printf("[ERROR] failed to exchange data1: %v\n", err)
	}
	if err := acks.SendData(client, server, data2); err != nil {
		fmt.Printf("[ERROR] failed to exchange data2: %v\n", err)
	}

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

	if _, err := forward.Send(client, server, req); err != nil {
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
	post1 := client.Engine.Protocol.BuildData([]byte("post-migration-1"))
	post2 := client.Engine.Protocol.BuildData([]byte("post-migration-2"))

	if err := acks.SendData(client, server, post1); err != nil {
		fmt.Printf("[ERROR] failed to exchange post1: %v\n", err)
	}
	if err := acks.SendData(client, server, post2); err != nil {
		fmt.Printf("[ERROR] failed to exchange post2: %v\n", err)
	}

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