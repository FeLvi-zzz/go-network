package ipv4

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/ipv4/types"
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type Consumer struct {
	config          *Config
	fragmentBuilder fragmentBuilder
	sender          sender
	icmpConsumer    icmpConsumer
	udpConsumer     udpConsumer
}

type icmpConsumer interface {
	Consume(b []byte, dstAddr []byte) (payload.Payload, error)
}

type udpConsumer interface {
	Consume(b []byte, ph []byte) (payload.Payload, error)
}

// TODO: set framents lifetime 15 sec
// https://datatracker.ietf.org/doc/html/rfc791#autoid-10
type fragmentBuilder map[uint16]map[int][]byte

func (f fragmentBuilder) Add(id uint16, offset int, b []byte) {
	if f[id] == nil {
		f[id] = make(map[int][]byte)
	}
	f[id][offset] = b
}

func (f fragmentBuilder) Build(id uint16) ([]byte, error) {
	defer delete(f, id)

	l := 0
	for _, b := range f[id] {
		l += len(b)
	}

	res := make([]byte, 0, l)
	offset := 0
	for {
		b := f[id][offset]
		delete(f[id], offset)

		res = append(res, b...)

		if b == nil {
			if len(f[id]) > 0 {
				return nil, fmt.Errorf("broken fragments")
			}
			return res, nil
		}

		offset = len(res)
	}
}

func NewConsumer(config *Config, sender sender, icmpConsumer icmpConsumer, udpConsumer udpConsumer) *Consumer {
	return &Consumer{
		config:          config,
		fragmentBuilder: make(fragmentBuilder),
		sender:          sender,
		icmpConsumer:    icmpConsumer,
		udpConsumer:     udpConsumer,
	}
}

func (c *Consumer) Consume(b []byte) (payload.Payload, error) {
	v4p, rb, err := FromBytes(b)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}
	if v4p.DstAddr != types.Address(c.config.localPrtAddr) {
		// debug: ignote the other dest packet
		return nil, fmt.Errorf("%w", util.ErrIgnorablePacket)
	}

	c.fragmentBuilder.Add(v4p.Identification, int(v4p.FlagmentOffset)*8, rb)

	if v4p.Flags.IsMF() {
		v4p.Payload = payload.NewFragmentPayload(len(rb))
	} else {
		nrb, err := c.fragmentBuilder.Build(v4p.Identification)
		if err != nil {
			return payload.NewFragmentPayload(len(rb)), err
		}

		switch v4p.Protocol {
		case types.Protocol_ICMP:
			icp, err := c.icmpConsumer.Consume(nrb, v4p.SrcAddr[:])
			if err != nil {
				return payload.NewUnknownPayload(nrb), err
			}
			v4p.Payload = icp
		case types.Protocol_UDP:
			up, err := c.udpConsumer.Consume(nrb, v4p.genUdpPseudoHeader(len(nrb)))
			if err != nil {
				return payload.NewUnknownPayload(nrb), err
			}
			v4p.Payload = up
		default:
			v4p.Payload = payload.NewUnknownPayload(nrb)
			// debug: ignore except for icmp
			return v4p, fmt.Errorf("%w", util.ErrIgnorablePacket)
		}
	}

	return v4p, nil
}
