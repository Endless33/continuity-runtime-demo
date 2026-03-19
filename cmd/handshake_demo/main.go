package main

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("HANDSHAKE DEMO (LOSSY EXCHANGE + KEEPALIVE + CLOSE)")
	fmt.Println()

	client := runtime.NewNode("client")
	server := runtime.NewNode("server")
	ex := runtime.NewLossyExchange("handshake")

	initPkt := client.StartHandshake()

	fmt.Println("\n=== SERVER RECEIVES INIT ===")
	resp, err := ex.Send(client, server, initPkt)
	if err != nil {
		fmt.Printf("[ERROR] server failed to handle init: %v\n", err)
		return
	}

	if resp == nil {
		fmt.Println("[WARN] init packet dropped or no init ack produced")
		return
	}

	fmt.Println("\n=== CLIENT RECEIVES INIT ACK ===")
	_, err = ex.Send(server, client, *resp)
	if err != nil {
		fmt.Printf("[ERROR] client failed to handle init ack: %v\n", err)
		return
	}

	fmt.Println("\n=== KEEPALIVE EXCHANGE ===")
	k1 := client.Engine.SendKeepalive()
	if k1.Type != "" {
		if _, err := ex.Send(client, server, k1); err != nil {
			fmt.Printf("[ERROR] client keepalive failed: %v\n", err)
		}
	}

	k2 := server.Engine.SendKeepalive()
	if k2.Type != "" {
		if _, err := ex.Send(server, client, k2); err != nil {
			fmt.Printf("[ERROR] server keepalive failed: %v\n", err)
		}
	}

	fmt.Println("\n=== CLOSE SESSION ===")
	closePkt := client.Engine.CloseSession(protocol.CloseReasonNormal)
	if _, err := ex.Send(client, server, closePkt); err != nil {
		fmt.Printf("[ERROR] close delivery failed: %v\n", err)
	}

	fmt.Println("\n=== RESULT ===")
	fmt.Printf("[CLIENT] session=%s state=%s epoch=%d path=%s\n",
		client.Engine.Protocol.SessionID,
		client.Engine.Protocol.State,
		client.Engine.Protocol.Epoch,
		client.Engine.Protocol.ActivePath,
	)

	fmt.Printf("[SERVER] session=%s state=%s epoch=%d path=%s\n",
		server.Engine.Protocol.SessionID,
		server.Engine.Protocol.State,
		server.Engine.Protocol.Epoch,
		server.Engine.Protocol.ActivePath,
	)

	fmt.Println("\n=== CLIENT TIMELINE ===")
	client.Engine.Runtime.Trace.PrintTimeline()

	fmt.Println("\n=== SERVER TIMELINE ===")
	server.Engine.Runtime.Trace.PrintTimeline()
}