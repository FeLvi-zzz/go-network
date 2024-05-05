package udp

import (
	"github.com/FeLvi-zzz/go-network/payload"
)

type Consumer struct {
	config *Config
	sender sender
}

type sender interface {
	UDPSend(targetAddr []byte, payload payload.Payload) error
}

func NewConsumer(config *Config, sender sender) *Consumer {
	return &Consumer{
		config: config,
		sender: sender,
	}
}

func (c *Consumer) Consume(b []byte, ph []byte) (payload.Payload, error) {
	up, err := FromBytes(b, ph)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	// TODO: pass to conn manager

	return up, nil
}
