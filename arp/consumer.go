package arp

import (
	"bytes"

	"github.com/FeLvi-zzz/go-network/payload"
)

type Consumer struct {
	config *Config
	sender sender
}

func NewConsumer(config *Config, sender sender) *Consumer {
	return &Consumer{
		config: config,
		sender: sender,
	}
}

func (c *Consumer) Consume(b []byte) (payload.Payload, error) {
	pap, err := FromBytes(b)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	if !(pap.Op == ArpOp_REQUEST && bytes.Equal(pap.Tpa, c.config.localPrtAddr)) {
		return pap, nil
	}

	nap := NewPayload(
		pap.Pro,
		ArpOp_RESPONSE,
		c.config.localHrdAddr,
		c.config.localPrtAddr,
		pap.Sha,
		pap.Tpa,
	)

	if err := c.sender.ArpSend(pap.Sha, nap); err != nil {
		return pap, err
	}

	return pap, nil
}
