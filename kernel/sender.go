package kernel

import (
	"fmt"
	"syscall"

	"github.com/FeLvi-zzz/go-network/payload"
)

type Sender struct {
	config *Config
}

func NewSender(config *Config) *Sender {
	return &Sender{
		config: config,
	}
}

func (s *Sender) Send(payload payload.Payload) error {
	fmt.Printf("\n-- send packet --\n")
	payload.Inspect()

	if err := syscall.Sendto(s.config.fd, payload.Bytes(), syscall.MSG_CONFIRM, &syscall.SockaddrLinklayer{
		Ifindex: s.config.ifIndex,
	}); err != nil {
		return err
	}

	return nil
}
