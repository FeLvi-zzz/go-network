package ethernet

import (
	"github.com/FeLvi-zzz/go-network/ethernet/types"
	"github.com/FeLvi-zzz/go-network/payload"
)

type Sender struct {
	config *Config
	sender sender
}

func NewSender(config *Config, sender sender) *Sender {
	return &Sender{
		config: config,
		sender: sender,
	}
}

func (s *Sender) ArpSend(targetHwAddr []byte, payload payload.Payload) error {
	if targetHwAddr == nil {
		targetHwAddr = s.config.BroadcastHrdAddr
	}

	eh := NewEthernetFrame(
		targetHwAddr,
		s.config.LocalHrdAddr,
		types.EtherType_ARP,
		payload,
	)
	return s.sender.Send(eh)
}
