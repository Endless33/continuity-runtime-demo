package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"continuity-runtime-demo/internal/session"
	"continuity-runtime-demo/internal/transport"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage:")
		fmt.Println("  server: go run ./cmd/udp_handshake/main.go server :9000")
		fmt.Println("  client: go run ./cmd/udp_handshake/main.go client :9001 127.0.0.1:9000")
		return
	}

	role := os.Args[1]
	listen := os.Args[2]

	node, err := transport.NewUDPNode(listen)
	if err != nil {
		panic(err)
	}

	switch role {
	case "server":
		runServer(node, listen)
	case "client":
		if len(os.Args) < 4 {
			fmt.Println("client mode requires remote address")
			return
		}
		remote := os.Args[3]
		runClient(node, listen, remote)
	default:
		fmt.Println("unknown role:", role)
	}
}

func runServer(node *transport.UDPNode, listen string) {
	fmt.Println("server listening on", listen)

	sessions := map[string]*session.Session{}

	node.Receive(func(data []byte, from *net.UDPAddr) {
		pkt, err := transport.DecodePacket(data)
		if err != nil {
			fmt.Println("decode error:", err)
			return
		}

		fmt.Printf("server recv type=%s session=%s epoch=%d from=%s\n",
			pkt.Type, pkt.SessionID, pkt.Epoch, from.String())

		sess, ok := sessions[pkt.SessionID]
		if !ok {
			sess = session.NewSession(pkt.SessionID)
			sessions[pkt.SessionID] = sess
		}

		sess.UpdatePath(from.String())

		switch pkt.Type {
		case transport.PacketTypeInit:
			sess.Establish()

			reply := transport.Packet{
				Type:      transport.PacketTypeInitAck,
				SessionID: sess.ID,
				Epoch:     sess.Epoch,
				Payload:   "ack",
			}

			raw, err := transport.EncodePacket(reply)
			if err != nil {
				fmt.Println("encode error:", err)
				return
			}

			if err := node.Send(from.String(), raw); err != nil {
				fmt.Println("send error:", err)
				return
			}

			fmt.Println("server sent init_ack to", from.String())

		case transport.PacketTypeData:
			if sess.State != session.StateEstablished {
				fmt.Println("server rejected data: session not established")
				return
			}
			fmt.Println("server data payload:", pkt.Payload)
		}
	})
}

func runClient(node *transport.UDPNode, listen, remote string) {
	fmt.Println("client listening on", listen)
	fmt.Println("client target:", remote)

	sess := session.NewSession("session-1")

	go node.Receive(func(data []byte, from *net.UDPAddr) {
		pkt, err := transport.DecodePacket(data)
		if err != nil {
			fmt.Println("decode error:", err)
			return
		}

		fmt.Printf("client recv type=%s session=%s epoch=%d from=%s\n",
			pkt.Type, pkt.SessionID, pkt.Epoch, from.String())

		if pkt.SessionID != sess.ID {
			fmt.Println("client ignored packet: wrong session")
			return
		}

		sess.UpdatePath(from.String())

		if pkt.Type == transport.PacketTypeInitAck {
			sess.Establish()

			dataPkt := transport.Packet{
				Type:      transport.PacketTypeData,
				SessionID: sess.ID,
				Epoch:     sess.Epoch,
				Payload:   "hello after handshake",
			}

			raw, err := transport.EncodePacket(dataPkt)
			if err != nil {
				fmt.Println("encode error:", err)
				return
			}

			if err := node.Send(remote, raw); err != nil {
				fmt.Println("send error:", err)
				return
			}

			fmt.Println("client sent data packet")
		}
	})

	initPkt := transport.Packet{
		Type:      transport.PacketTypeInit,
		SessionID: sess.ID,
		Epoch:     sess.Epoch,
		Payload:   "hello",
	}

	raw, err := transport.EncodePacket(initPkt)
	if err != nil {
		panic(err)
	}

	time.Sleep(150 * time.Millisecond)

	if err := node.Send(remote, raw); err != nil {
		panic(err)
	}

	fmt.Println("client sent init")

	select {}
}