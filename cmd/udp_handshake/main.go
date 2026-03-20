package main

import (
	"fmt"
	"net"
	"os"

	"continuity-runtime-demo/internal/session"
	"continuity-runtime-demo/internal/transport"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run cmd/udp_handshake/main.go :9000")
		return
	}

	listen := os.Args[1]

	node, err := transport.NewUDPNode(listen)
	if err != nil {
		panic(err)
	}

	sess := session.NewSession("session-1")

	fmt.Println("🚀 UDP node listening on", listen)

	node.Receive(func(data []byte, from *net.UDPAddr) {
		addr := from.String()

		fmt.Println("📦 recv from", addr, ":", string(data))

		// 🔥 ключевая идея: session отслеживает path
		sess.UpdatePath(addr)

		// имитация authority / epoch логики
		fmt.Println("🧠 session:", sess.ID, "epoch:", sess.Epoch)

		// отправка ответа
		err := node.Send(addr, []byte("ack:"+sess.ID))
		if err != nil {
			fmt.Println("send error:", err)
		}
	})
}