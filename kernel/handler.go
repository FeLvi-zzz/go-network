package kernel

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type Handler struct {
	config      *Config
	logger      *Logger
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
		recvTime := time.Now()

		b := make([]byte, 1518) // MTU 1500 + Ethernet header 14 + Ethernet FCS 4
		_, _, err := syscall.Recvfrom(h.config.fd, b, 0)
		if err != nil {
			return err
		}

		li := logger.Reserve()

		p, err := h.ethConsumer.Consume(b)
		if errors.Is(err, payload.ErrUnknownPayload) {
			logger.Log(li, func() {})
			continue
		}
		if errors.Is(err, util.ErrIgnorablePacket) {
			logger.Log(li, func() {})
			continue
		}

		logger.Log(li, func() {
			fmt.Printf("\n-- recv packet --\n%s\n", recvTime.String())
			p.Inspect()
		})

		if err != nil {
			return err
		}
	}
}
