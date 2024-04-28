package arp

type Config struct {
	localHrdAddr []byte
	localPrtAddr []byte
}

func NewConfig(localHrdAddr []byte, localPrtAddr []byte) *Config {
	return &Config{
		localHrdAddr: localHrdAddr,
		localPrtAddr: localPrtAddr,
	}
}
