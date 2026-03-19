package main

import (
	"fmt"

	"continuity-runtime-demo/internal/runtime"
)

func main() {
	fmt.Println("HANDSHAKE DEMO")
	fmt.Println()

	client := runtime.NewNode("client")
	server := runtime.NewNode("server")
	ex := runtime.NewExchange("handshake")

	initPkt := client.StartHandshake()

	fmt.Println("\n=== SERVER RECEIVES INIT ===")
	resp, err := ex.Send(client, server, initPkt)
	if err != nil {
		fmt.Printf("[ERROR] server failed to handle init: %v\n", err)
		return
	}

	if resp == nil {
		fmt.Println("[ERROR] server did not produce init ack")
		return
	}

	fmt.Println("\n=== CLIENT RECEIVES INIT ACK ===")
	_, err = ex.Send(server, client, *resp)
	if err != nil {
		fmt.Printf("[ERROR] client failed to handle init ack: %v\n", err)
		return
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