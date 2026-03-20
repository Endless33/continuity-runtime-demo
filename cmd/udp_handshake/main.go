package main

import (
	"fmt"
	"os"

	"yourmodule/internal/session"
	"yourmodule/internal/transport"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go [listen_addr]")
		return
	}

	listen := os.Args[1]

	node, err := transport.NewUDPNode(listen)
	if err != nil {
		panic(err)
	}

	sess := session.NewSession("session-1")

	fmt.Println("listening on", listen)

	node.Receive(func(data []byte, fromAddr interface{}) {
		addr := fromAddr.(*net.UDPAddr).String()

		fmt.Println("📦 recv from", addr, ":", string(data))

		// 🔥 ключевая логика
		sess.UpdatePath(addr)

		// отправляем ответ
		node.Send(addr, []byte("ack:"+sess.ID))
	})
}