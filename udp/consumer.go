package udp

import (
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/udp/types"
)

type Consumer struct {
	config *Config
	sender sender
}

type sender interface {
	UDPSend(targetAddr []byte, payload payload.Payload) error
	UDPPseudoHeader(srcAddr []byte, dstAddr []byte, datalen int) []byte
}

func NewConsumer(config *Config, sender sender) *Consumer {
	return &Consumer{
		config: config,
		sender: sender,
	}
}

func (c *Consumer) Consume(b []byte, ph []byte, srcAddr []byte, dstAddr []byte) (payload.Payload, error) {
	up, err := FromBytes(b, ph)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	la := types.Address{
		IP:   dstAddr,
		Port: up.DstPort,
	}
	ra := types.Address{
		IP:   srcAddr,
		Port: up.SrcPort,
	}

	listener := LookupListener(la)
	if listener == nil {
		return up, nil
	}

	conn, ok := listener.conns[ra.String()]
	if ok {
		conn.data = up.Data.Bytes()
	} else {
		listener.conns[ra.String()] = &Conn{
			laddr:  la,
			raddr:  ra,
			data:   up.Data.Bytes(),
			sender: c.sender,
			cleanup: func() error {
				delete(listener.conns, ra.String())
				return nil
			},
		}
	}

	listener.receiver <- ra

	return up, nil
}
