package ipv4

import (
	"github.com/FeLvi-zzz/go-network/ipv4/types"
	"github.com/FeLvi-zzz/go-network/payload"
)

type Sender struct {
	config  *Config
	sender  sender
	ar      addressResolver
	idCount uint16
}

type sender interface {
	IPv4Send(targetAddr []byte, payload payload.Payload) error
}

type addressResolver interface {
	Resolve(targetAddr []byte) ([]byte, error)
}

func NewSender(config *Config, sender sender, ar addressResolver) *Sender {
	return &Sender{
		config: config,
		sender: sender,
		ar:     ar,
	}
}

func (s *Sender) send(targetAddr []byte, protocol types.Protocol, payload payload.Payload) error {
	defer func() {
		s.idCount++
	}()

	p := NewIPv4Packet(
		types.Address(targetAddr),
		types.Address(s.config.localPrtAddr),
		s.idCount,
		protocol,
		payload,
	)

	ntaddr := s.config.rt.Resolve(types.Address(targetAddr))

	addr, err := s.ar.Resolve(ntaddr[:])
	if err != nil {
		return err
	}

	if err := s.sender.IPv4Send(addr, p); err != nil {
		return err
	}

	return nil
}

func (s *Sender) ICMPSend(targetAddr []byte, payload payload.Payload) error {
	return s.send(targetAddr, types.Protocol_ICMP, payload)
}

func (s *Sender) UDPSend(targetAddr []byte, payload payload.Payload) error {
	return s.send(targetAddr, types.Protocol_UDP, payload)
}

func (s *Sender) UDPPseudoHeader(srcAddr []byte, dstAddr []byte, datalen int) []byte {
	b := make([]byte, 0, 12)
	b = append(b, srcAddr...)
	b = append(b, dstAddr...)
	b = append(b, 0)
	b = append(b, byte(types.Protocol_UDP))
	b = append(b, byte(datalen>>8))
	b = append(b, byte(datalen))

	return b
}

func (s *Sender) TCPSend(targetAddr []byte, payload payload.Payload) error {
	return s.send(targetAddr, types.Protocol_TCP, payload)
}

func (s *Sender) TCPPseudoHeader(srcAddr []byte, dstAddr []byte, datalen int) []byte {
	b := make([]byte, 0, 12)
	b = append(b, srcAddr...)
	b = append(b, dstAddr...)
	b = append(b, 0)
	b = append(b, byte(types.Protocol_TCP))
	b = append(b, byte(datalen>>8))
	b = append(b, byte(datalen))

	return b
}
