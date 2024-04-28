package ethernet

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/ethernet/types"
	"github.com/FeLvi-zzz/go-network/payload"
)

type Consumer struct {
	config      *Config
	arpConsumer arpConsumer
	sender      sender
}

type arpConsumer interface {
	Consume([]byte) (payload.Payload, error)
}

type sender interface {
	Send(payload payload.Payload) error
}

func NewConsumer(config *Config, arpConsumer arpConsumer, sender sender) *Consumer {
	return &Consumer{
		config:      config,
		arpConsumer: arpConsumer,
		sender:      sender,
	}
}

func (c *Consumer) Consume(b []byte) (payload.Payload, error) {
	ef, rb, err := FromBytes(b)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	if ef.Ethertype == types.EtherType_ARP {
		ap, err := c.arpConsumer.Consume(rb)
		if err != nil {
			return payload.NewUnknownPayload(rb), err
		}

		ef.Payload = ap

		return ef, nil
	}

	return payload.NewUnknownPayload(b), fmt.Errorf("ethernet: %w", payload.ErrUnknownPayload)
}
