package udp

import "github.com/FeLvi-zzz/go-network/udp/types"

type Service struct {
	sender sender
}

func NewService(sender sender) *Service {
	return &Service{
		sender: sender,
	}
}

func (s *Service) Listen(addr []byte, port uint16) *Listener {
	laddr := types.Address{
		IP:   addr,
		Port: port,
	}

	l := &Listener{
		receiver: make(chan types.Address, 1),
		conns:    make(map[string]*Conn),
		sender:   s.sender,
		laddr:    laddr,
	}

	listenerMap[laddr.String()] = l

	return l
}

func (s *Service) Dial(raddr []byte, rport uint16, laddr []byte, lport uint16) *Conn {
	la := types.Address{
		IP:   laddr,
		Port: lport,
	}
	ra := types.Address{
		IP:   raddr,
		Port: rport,
	}

	l := s.Listen(laddr, lport)
	c := &Conn{
		laddr:    la,
		raddr:    ra,
		sender:   s.sender,
		dataChan: make(chan []byte, 1),
		cleanup: func() error {
			delete(l.conns, ra.String())
			delete(listenerMap, la.String())
			return nil
		},
	}

	l.conns[ra.String()] = c

	return c
}
