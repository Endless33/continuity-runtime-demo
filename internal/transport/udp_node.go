package transport

import (
	"fmt"
	"net"
)

type UDPNode struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func NewUDPNode(listenAddr string) (*UDPNode, error) {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &UDPNode{
		conn: conn,
		addr: addr,
	}, nil
}

func (n *UDPNode) Send(to string, data []byte) error {
	addr, err := net.ResolveUDPAddr("udp", to)
	if err != nil {
		return err
	}

	_, err = n.conn.WriteToUDP(data, addr)
	return err
}

func (n *UDPNode) Receive(handler func(data []byte, from *net.UDPAddr)) {
	buf := make([]byte, 2048)

	for {
		nRead, addr, err := n.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		handler(buf[:nRead], addr)
	}
}