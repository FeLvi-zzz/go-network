package ipv4

type Config struct {
	localPrtAddr []byte
}

func NewConfig(localPrtAddr []byte) *Config {
	return &Config{
		localPrtAddr: localPrtAddr,
	}
}
