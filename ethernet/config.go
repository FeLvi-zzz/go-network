package ethernet

type Config struct {
	LocalHrdAddr     []byte
	BroadcastHrdAddr []byte
}

func NewConfig(localHrdAddr []byte, broadcastHrdAddr []byte) *Config {
	return &Config{
		LocalHrdAddr:     localHrdAddr,
		BroadcastHrdAddr: broadcastHrdAddr,
	}
}
