package ipv4

import (
	"github.com/FeLvi-zzz/go-network/payload"
)

type Consumer struct {
	config *Config
	// sender sender
}

func NewConsumer(config *Config) *Consumer {
	return &Consumer{
		config: config,
	}
}

func (c *Consumer) Consume(b []byte) (payload.Payload, error) {
	v4p, rb, err := FromBytes(b)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	// TODO
	v4p.Payload = payload.NewUnknownPayload(rb)

	return v4p, nil
}
