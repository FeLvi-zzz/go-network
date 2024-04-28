package arp

import (
	ethtypes "github.com/FeLvi-zzz/go-network/ethernet/types"
)

type Handler struct {
	config *Config
	sender sender
}

func NewHandler(config *Config, sender sender) *Handler {
	return &Handler{
		config: config,
		sender: sender,
	}
}

func (h *Handler) Startup(dstPrtAddr []byte) error {
	ap := NewPayload(
		ethtypes.EtherType_IPv4,
		ArpOp_REQUEST,
		h.config.localHrdAddr,
		h.config.localPrtAddr,
		[]byte{0, 0, 0, 0, 0, 0},
		dstPrtAddr,
	)

	h.sender.ArpSend(nil, ap)

	return nil
}
