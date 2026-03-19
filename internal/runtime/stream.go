package runtime

type Stream struct {
	Runtime *Runtime
	Network *NetworkSimulator
}

func NewStream(r *Runtime, n *NetworkSimulator) *Stream {
	return &Stream{
		Runtime: r,
		Network: n,
	}
}

func (s *Stream) Send(n int) {
	for i := 0; i < n; i++ {
		s.Runtime.PacketID++

		packetID := s.Runtime.PacketID

		s.Network.Transmit(packetID, s.Runtime.Current)
	}
}