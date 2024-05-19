package tcp

import "github.com/FeLvi-zzz/go-network/tcp/types"

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
		conns:    connMap{},
		sender:   s.sender,
		laddr:    laddr,
	}

	globalListenerMap.Store(laddr.String(), l)

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
		dataChan: make(chan []byte, 100),
		cleanup: func() error {
			l.conns.Delete(ra.String())
			globalListenerMap.Delete(la.String())
			return nil
		},
	}

	// l.conns[ra.String()] = c

	return c
}
