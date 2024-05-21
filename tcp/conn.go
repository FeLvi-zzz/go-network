package tcp

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

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
	data := <-c.dataChan

	if c.state == types.State_CLOSE_WAIT || c.state == types.State_CLOSED {
		return copy(p, data), fmt.Errorf("network closed")
	}

	return copy(p, data), io.EOF
}

func (c *Conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.state = types.State_LAST_ACK
	seg := NewSegment(c.laddr.Port, c.raddr.Port, c.sndnxt, c.rcvnxt, types.Flags_ACK|types.Flags_FIN, payload.NewDataPayload(nil))
	c.sndnxt += 1
	if err := c.send(seg); err != nil {
		return err
	}

	return nil
}

func (c *Conn) Write(p []byte) (n int, err error) {
	for {
		if c.state == types.State_ESTAB {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
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
		if ts.Flags&types.Flags_ACK != 0 {
			if ts.AckNum <= c.snduna || ts.AckNum > c.sndnxt {
				ns := NewSegment(c.laddr.Port, c.raddr.Port, ts.AckNum, 0, types.Flags_RST, payload.NewDataPayload(nil))
				return c.send(ns)
			}
		}

		if ts.Flags&types.Flags_RST != 0 {
			return c.cleanup()
		}

		if ts.Flags&types.Flags_SYN != 0 {
			if ts.Flags&types.Flags_ACK != 0 {
				c.rcvnxt = ts.SeqNum + 1
				c.snduna = ts.AckNum
				c.state = types.State_ESTAB

				ns := NewSegment(c.laddr.Port, c.raddr.Port, c.sndnxt, c.rcvnxt, types.Flags_ACK, payload.NewDataPayload(nil))
				return c.send(ns)
			} else {
				c.state = types.State_SYN_RCVD
				c.sndwnd = ts.Window
				c.sndwl1 = ts.SeqNum
				c.sndwl2 = ts.AckNum

				ns := NewSegment(c.laddr.Port, c.raddr.Port, c.sndnxt, c.rcvnxt, types.Flags_ACK, payload.NewDataPayload(nil))
				return c.send(ns)
			}
		}
		return nil
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
			c.listener.receiver <- c.raddr
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
		}(c)

		if len(ts.Data.Bytes()) == 0 && ts.Flags&types.Flags_FIN == 0 {
			return nil
		}

		c.rcvnxt = ts.SeqNum + uint32(len(ts.Data.Bytes()))

		if ts.Flags&types.Flags_FIN != 0 {
			c.state = types.State_CLOSE_WAIT
			c.rcvnxt += 1
			close(c.dataChan)
		}

		ns := NewSegment(
			c.laddr.Port,
			c.raddr.Port,
			c.sndnxt,
			c.rcvnxt,
			types.Flags_ACK,
			payload.NewDataPayload(nil),
		)

		return c.send(ns)
	case types.State_LAST_ACK:
		c.state = types.State_CLOSED
		return c.cleanup()
	}

	return nil
}

func (c *Conn) send(ns *Segment) error {
	ns.PseudoHeader = c.sender.TCPPseudoHeader(c.laddr.IP, c.raddr.IP, len(ns.Bytes()))
	ns.Checksum = ns.CalcCheckSum()
	return c.sender.TCPSend(c.raddr.IP, ns)
}
