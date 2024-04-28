package ethernet

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/ethernet/types"
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type Consumer struct {
	config       *Config
	arpConsumer  arpConsumer
	ipv4Consumer ipv4Consumer
	sender       sender
}

type arpConsumer interface {
	Consume([]byte) (payload.Payload, error)
}

type ipv4Consumer interface {
	Consume([]byte) (payload.Payload, error)
}

type sender interface {
	Send(payload payload.Payload) error
}

func NewConsumer(config *Config, arpConsumer arpConsumer, ipv4Consumer ipv4Consumer, sender sender) *Consumer {
	return &Consumer{
		config:       config,
		arpConsumer:  arpConsumer,
		ipv4Consumer: ipv4Consumer,
		sender:       sender,
	}
}

func (c *Consumer) Consume(b []byte) (payload.Payload, error) {
	ef, rb, err := FromBytes(b)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	// debug: ignore loopback
	if ef.DstMacAddr == [6]byte{} && ef.SrcMacAddr == [6]byte{} {
		return nil, fmt.Errorf("%w", util.ErrIgnorablePacket)
	}

	switch ef.Ethertype {
	case types.EtherType_ARP:
		ap, err := c.arpConsumer.Consume(rb)
		if err != nil {
			return payload.NewUnknownPayload(rb), err
		}

		ef.Payload = ap

		return ef, nil
	case types.EtherType_IPv4:
		v4p, err := c.ipv4Consumer.Consume(rb)
		if err != nil {
			return payload.NewUnknownPayload(rb), err
		}

		ef.Payload = v4p

		return ef, nil
	}

	return payload.NewUnknownPayload(b), fmt.Errorf("ethernet: %w", payload.ErrUnknownPayload)
}
