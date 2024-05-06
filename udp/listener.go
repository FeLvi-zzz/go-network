package udp

import "github.com/FeLvi-zzz/go-network/udp/types"

var (
	listenerMap = make(map[string]*Listener)
)

func LookupListener(addr types.Address) *Listener {
	return listenerMap[addr.String()]
}

type Listener struct {
	receiver chan types.Address
	conns    map[string]*Conn
	sender   sender
	laddr    types.Address
}

func (l *Listener) Accept() *Conn {
	ra := <-l.receiver
	return l.conns[ra.String()]
}
