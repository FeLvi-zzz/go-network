package tcp

import (
	"fmt"
	"io"
	"math/rand"
	"sync"

	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/tcp/types"
)

type Conn struct {
	mu           sync.Mutex
	listener     *Listener
	state        types.State
	isActiveOpen bool
	laddr        types.Address
	raddr        types.Address
	sndnxt       uint32
	snduna       uint32
	sndwnd       uint16
	sndwl1       uint32
	sndwl2       uint32
	rcvnxt       uint32
	dataChan     chan []byte
	sender       sender
	cleanup      func() error
}

func (c *Conn) Read(p []byte) (n int, err error) {
	return copy(p, <-c.dataChan), io.EOF
}

func (c *Conn) Close() error {
	return c.cleanup()
}

func (c *Conn) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data := payload.NewDataPayload(p)
	seg := NewSegment(c.laddr.Port, c.raddr.Port, c.sndnxt, c.rcvnxt, types.Flags_ACK, data)
	c.sndnxt += uint32(len(data.Bytes()))
	if err := c.send(seg); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (c *Conn) consume(ts *Segment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("state: %d\n", c.state)
	switch c.state {
	case types.State_LISTEN:
		if ts.Flags&types.Flags_RST != 0 {
			return nil
		}
		if ts.Flags&types.Flags_ACK != 0 {
			ns := NewSegment(c.laddr.Port, c.raddr.Port, ts.AckNum, 0, types.Flags_RST, payload.NewDataPayload(nil))
			return c.send(ns)
		}
		if ts.Flags&types.Flags_SYN != 0 {
			iss := rand.Uint32()
			c.sndnxt = iss
			c.rcvnxt = ts.SeqNum + 1
			ns := NewSegment(c.laddr.Port, c.raddr.Port, c.sndnxt, c.rcvnxt, types.Flags_SYN|types.Flags_ACK, payload.NewDataPayload(nil))
			c.snduna = iss
			c.sndnxt += 1
			c.state = types.State_SYN_RCVD

			return c.send(ns)
		}
	case types.State_SYN_SENT:
		return fmt.Errorf("not implemented")
	}

	switch c.state {
	case types.State_SYN_RCVD:
		if ts.Flags&types.Flags_RST != 0 {
			if c.isActiveOpen {
				return c.cleanup()
			} else {
				c.state = types.State_LISTEN
				return nil
			}
		}

		if ts.Flags&types.Flags_SYN != 0 {
			if !c.isActiveOpen {
				c.state = types.State_LISTEN
			}
			return nil
		}
	}

	if ts.Flags&types.Flags_ACK == 0 {
		return nil
	}

	switch c.state {
	case types.State_SYN_RCVD:
		if c.snduna < ts.AckNum && ts.AckNum >= c.sndnxt {
			c.state = types.State_ESTAB
			c.sndwl1 = ts.SeqNum
			c.sndwl2 = ts.AckNum
		} else {
			ns := NewSegment(c.laddr.Port, c.raddr.Port, ts.AckNum, 0, types.Flags_RST, payload.NewDataPayload(nil))
			return c.send(ns)
		}
	case types.State_ESTAB:
		if c.snduna <= ts.AckNum && ts.AckNum <= c.sndnxt {
			c.snduna = ts.AckNum

			if c.sndwl1 < ts.SeqNum || (c.sndwl1 == ts.SeqNum && c.sndwl2 <= ts.AckNum) {
				c.sndwnd = ts.Window
				c.sndwl1 = ts.SeqNum
				c.sndwl2 = ts.AckNum
			}
		}

		defer func(c *Conn) {
			if len(ts.Data.Bytes()) == 0 {
				return
			}
			c.dataChan <- ts.Data.Bytes()
			c.listener.receiver <- c.raddr
		}(c)

		if len(ts.Data.Bytes()) == 0 {
			return nil
		}

		c.rcvnxt = ts.SeqNum + uint32(len(ts.Data.Bytes()))

		ns := NewSegment(
			c.laddr.Port,
			c.raddr.Port,
			c.sndnxt,
			c.rcvnxt,
			types.Flags_ACK,
			payload.NewDataPayload(nil),
		)
		return c.send(ns)
	}

	return nil
}

func (c *Conn) send(ns *Segment) error {
	ns.PseudoHeader = c.sender.TCPPseudoHeader(c.laddr.IP, c.raddr.IP, len(ns.Bytes()))
	ns.Checksum = ns.CalcCheckSum()
	return c.sender.TCPSend(c.raddr.IP, ns)
}
