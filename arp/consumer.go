package arp

import (
	"bytes"

	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
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

	if !bytes.Equal(pap.Tpa, c.config.localPrtAddr) {
		return pap, util.ErrIgnorablePacket
	}

	if pap.Op == ArpOp_RESPONSE {
		globalArpTable[[4]byte(pap.Spa)] = [6]byte(pap.Sha)
		conn.res <- pap.Sha
		return pap, nil
	}

	nap := NewPayload(
		pap.Pro,
		ArpOp_RESPONSE,
		c.config.localHrdAddr,
		c.config.localPrtAddr,
		pap.Sha,
		pap.Spa,
	)

	if err := c.sender.ArpSend(pap.Sha, nap); err != nil {
		return pap, err
	}

	return pap, nil
}
