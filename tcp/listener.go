package tcp

import (
	"sync"

	"github.com/FeLvi-zzz/go-network/tcp/types"
)

var (
	globalListenerMap = listenerMap{}
)

func LookupListener(addr types.Address) *Listener {
	l, _ := globalListenerMap.Load(addr.String())
	return l
}

type listenerMap struct {
	m sync.Map
}

func (m *listenerMap) Load(key string) (*Listener, bool) {
	value, ok := m.m.Load(key)
	return value.(*Listener), ok
}

func (m *listenerMap) LoadOrStore(key string, value *Listener) (*Listener, bool) {
	actual, loaded := m.m.LoadOrStore(key, value)
	return actual.(*Listener), loaded
}

func (m *listenerMap) Store(key string, value *Listener) {
	m.m.Store(key, value)
}

func (m *listenerMap) Delete(key string) {
	m.m.Delete(key)
}

type connMap struct {
	m sync.Map
}

func (m *connMap) Load(key string) (*Conn, bool) {
	value, ok := m.m.Load(key)
	return value.(*Conn), ok
}

func (m *connMap) Store(key string, value *Conn) {
	m.m.Store(key, value)
}

func (m *connMap) LoadOrStore(key string, value *Conn) (*Conn, bool) {
	actual, loaded := m.m.LoadOrStore(key, value)
	return actual.(*Conn), loaded
}

func (m *connMap) Delete(key string) {
	m.m.Delete(key)
}

type Listener struct {
	receiver chan types.Address
	conns    connMap
	sender   sender
	laddr    types.Address
}

func (l *Listener) Accept() *Conn {
	ra := <-l.receiver
	c, _ := l.conns.Load(ra.String())
	return c
}

func (l *Listener) consume(ts *Segment, ra types.Address) error {
	conn, _ := l.conns.LoadOrStore(ra.String(), &Conn{
		listener:     l,
		state:        types.State_LISTEN,
		isActiveOpen: false,
		laddr:        l.laddr,
		raddr:        ra,
		dataChan:     make(chan []byte, 100),
		sender:       l.sender,
		cleanup: func() error {
			l.conns.Delete(ra.String())
			return nil
		},
	})
	if err := conn.consume(ts); err != nil {
		return err
	}
	return nil
}
