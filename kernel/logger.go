package kernel

import (
	"sync"
	"sync/atomic"
)

type Logger struct {
	n       atomic.Int64
	datamap sync.Map // int64 -> func()
	chanmap sync.Map // int64 -> chan struct{}
}

var (
	logger *Logger
)

func (l *Logger) Reserve() int64 {
	return l.n.Add(1)
}

func (l *Logger) Log(i int64, f func()) {
	l.datamap.Store(i, f)
	c, _ := l.chanmap.LoadOrStore(i, make(chan struct{}, 1))
	c.(chan struct{}) <- struct{}{}
}

func init() {
	logger = &Logger{}
	go func() {
		for i := int64(1); ; i++ {
			c, _ := logger.chanmap.LoadOrStore(i, make(chan struct{}, 1))
			<-c.(chan struct{})

			f, _ := logger.datamap.Load(i)
			f.(func())()
			logger.chanmap.Delete(i)
			logger.datamap.Delete(i)
		}
	}()
}
