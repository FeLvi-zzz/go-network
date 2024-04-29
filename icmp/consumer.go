package icmp

import (
	"github.com/FeLvi-zzz/go-network/icmp/messages"
	"github.com/FeLvi-zzz/go-network/icmp/types"
	"github.com/FeLvi-zzz/go-network/payload"
)

type Consumer struct {
	config *Config
	sender sender
}

type sender interface {
	ICMPSend(targetAddr []byte, payload payload.Payload) error
}

func NewConsumer(config *Config, sender sender) *Consumer {
	return &Consumer{
		config: config,
		sender: sender,
	}
}

func (c *Consumer) Consume(b []byte, dstAddr []byte) (payload.Payload, error) {
	ip, err := messages.FromBytes(b)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	switch ip.Type {
	case types.ICMPType_EchoReply:
	case types.ICMPType_Echo:
		p := messages.EchoPayloadFromBytes(ip.Payload.Bytes())
		ep := messages.NewEchoReplyMessage(p.Identifier, p.SequenceNumber, p.Data.Bytes())
		if err := c.sender.ICMPSend(dstAddr, ep); err != nil {
			return ip, err
		}
	}

	return ip, nil
}
