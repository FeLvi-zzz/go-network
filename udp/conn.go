package udp

import (
	"io"

	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/udp/types"
)

type Conn struct {
	laddr    types.Address
	raddr    types.Address
	dataChan chan []byte
	sender   sender
	cleanup  func() error
}

func (c *Conn) Read(p []byte) (n int, err error) {
	return copy(p, <-c.dataChan), io.EOF
}

func (c *Conn) Close() error {
	return c.cleanup()
}

func (c *Conn) Write(p []byte) (n int, err error) {
	d := NewDatagram(c.laddr.Port, c.raddr.Port, payload.NewDataPayload(p))

	d.PseudoHeader = c.sender.UDPPseudoHeader(c.laddr.IP, c.raddr.IP, 8+len(p))
	d.Checksum = d.CalcCheckSum()

	if err := c.sender.UDPSend(c.raddr.IP, d); err != nil {
		return 0, err
	}

	return len(p), nil
}
