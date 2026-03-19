package runtime

import (
	"fmt"

	"continuity-runtime-demo/internal/protocol"
)

type Node struct {
	Name      string
	Engine    *Engine
	Handshake *protocol.HandshakeEngine
	States    *protocol.StateMachine
}

func NewNode(name string) *Node {
	return &Node{
		Name:      name,
		Engine:    NewEngine(),
		Handshake: protocol.NewHandshakeEngine(),
		States:    protocol.NewStateMachine(),
	}
}

func (n *Node) StartHandshake() protocol.WirePacket {
	pkt := n.Handshake.BuildInit(n.Engine.Protocol)

	n.Engine.Runtime.Trace.Record("handshake_init_sent", "session init sent", map[string]interface{}{
		"node":       n.Name,
		"session_id": pkt.SessionID,
		"epoch":      pkt.Epoch,
		"path":       pkt.Path,
	})

	fmt.Printf("[%s] handshake init session=%s epoch=%d path=%s\n",
		n.Name,
		pkt.SessionID,
		pkt.Epoch,
		pkt.Path,
	)

	return pkt
}

func (n *Node) Receive(pkt protocol.WirePacket) (*protocol.WirePacket, error) {
	switch pkt.Type {
	case protocol.PacketTypeSessionInit:
		res := n.Handshake.HandleInit(n.Engine.Protocol, pkt)
		if !res.Accepted {
			return nil, fmt.Errorf(res.Reason)
		}

		if err := n.States.Apply(n.Engine.Protocol, protocol.SessionStateAttached); err != nil {
			return nil, err
		}

		n.Engine.Runtime.Trace.Record("handshake_init_received", "session init received", map[string]interface{}{
			"node":       n.Name,
			"session_id": pkt.SessionID,
			"epoch":      pkt.Epoch,
			"path":       pkt.Path,
		})

		fmt.Printf("[%s] handshake init accepted session=%s epoch=%d path=%s\n",
			n.Name,
			pkt.SessionID,
			pkt.Epoch,
			pkt.Path,
		)

		return &res.Packet, nil

	case protocol.PacketTypeSessionInitAck:
		res := n.Handshake.HandleInitAck(n.Engine.Protocol, pkt)
		if !res.Accepted {
			return nil, fmt.Errorf(res.Reason)
		}

		if err := n.States.Apply(n.Engine.Protocol, protocol.SessionStateAttached); err != nil {
			return nil, err
		}

		n.Engine.Runtime.Trace.Record("handshake_ack_received", "session init ack received", map[string]interface{}{
			"node":       n.Name,
			"session_id": pkt.SessionID,
			"epoch":      pkt.Epoch,
			"path":       pkt.Path,
		})

		fmt.Printf("[%s] handshake ack accepted session=%s epoch=%d path=%s\n",
			n.Name,
			pkt.SessionID,
			pkt.Epoch,
			pkt.Path,
		)

		return nil, nil

	default:
		return nil, n.Engine.Receive(pkt)
	}
}

func (n *Node) BeginRecovery() error {
	if err := n.States.Apply(n.Engine.Protocol, protocol.SessionStateRecovering); err != nil {
		return err
	}

	n.Engine.Runtime.Trace.Record("node_recovering", "node entered recovering state", map[string]interface{}{
		"node":  n.Name,
		"state": string(n.Engine.Protocol.State),
	})

	fmt.Printf("[%s] state -> %s\n", n.Name, n.Engine.Protocol.State)
	return nil
}

func (n *Node) EndRecovery() error {
	if err := n.States.Apply(n.Engine.Protocol, protocol.SessionStateAttached); err != nil {
		return err
	}

	n.Engine.Runtime.Trace.Record("node_attached", "node returned to attached state", map[string]interface{}{
		"node":  n.Name,
		"state": string(n.Engine.Protocol.State),
	})

	fmt.Printf("[%s] state -> %s\n", n.Name, n.Engine.Protocol.State)
	return nil
}