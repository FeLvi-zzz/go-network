package kernel

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type Handler struct {
	config      *Config
	ethConsumer ethernetConsumer
}

type ethernetConsumer interface {
	Consume(b []byte) (payload.Payload, error)
}

func NewHandler(config *Config, ethConsumer ethernetConsumer) *Handler {
	return &Handler{
		config:      config,
		ethConsumer: ethConsumer,
	}
}

func (h *Handler) Handle() error {
	for {
		b := make([]byte, 1518) // MTU 1500 + Ethernet header 14 + Ethernet FCS 4
		_, _, err := syscall.Recvfrom(h.config.fd, b, 0)
		if err != nil {
			return err
		}

		p, err := h.ethConsumer.Consume(b)
		if errors.Is(err, payload.ErrUnknownPayload) {
			continue
		}
		if errors.Is(err, util.ErrIgnorablePacket) {
			continue
		}

		fmt.Printf("\n-- recv packet --\n")
		p.Inspect()

		if err != nil {
			return err
		}
	}
}
