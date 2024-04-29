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
