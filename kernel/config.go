package kernel

import "sync"

type Config struct {
	fd      int
	ifIndex int
	mu      *sync.Mutex
}

func NewConfig(fd int, ifIndex int) *Config {
	return &Config{
		fd:      fd,
		ifIndex: ifIndex,
		mu:      &sync.Mutex{},
	}
}
