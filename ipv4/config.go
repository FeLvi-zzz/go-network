package ipv4

type Config struct {
	localPrtAddr []byte
	rt           RouteTable
}

func NewConfig(localPrtAddr []byte, rt RouteTable) *Config {
	return &Config{
		localPrtAddr: localPrtAddr,
		rt:           rt,
	}
}
