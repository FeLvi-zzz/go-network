package tcp

import (
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/tcp/types"
)

type Consumer struct {
	config *Config
	sender sender
}

type sender interface {
	TCPSend(targetAddr []byte, payload payload.Payload) error
	TCPPseudoHeader(srcAddr []byte, dstAddr []byte, datalen int) []byte
}

func NewConsumer(config *Config, sender sender) *Consumer {
	return &Consumer{
		config: config,
		sender: sender,
	}
}

func (c *Consumer) Consume(b []byte, ph []byte, srcAddr []byte, dstAddr []byte) (payload.Payload, error) {
	ts, err := FromBytes(b, ph)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	la := types.Address{
		IP:   dstAddr,
		Port: ts.DstPort,
	}
	ra := types.Address{
		IP:   srcAddr,
		Port: ts.SrcPort,
	}

	listener := LookupListener(la)
	if listener == nil {
		return ts, nil
	}

	if err := listener.consume(ts, ra); err != nil {
		return ts, err
	}

	return ts, nil
}
