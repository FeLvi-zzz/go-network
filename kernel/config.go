package kernel

type Config struct {
	fd      int
	ifIndex int
}

func NewConfig(fd int, ifIndex int) *Config {
	return &Config{
		fd:      fd,
		ifIndex: ifIndex,
	}
}
