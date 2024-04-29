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

func (h *Handler) Request(dstPrtAddr []byte) error {
	ap := NewPayload(
		ethtypes.EtherType_IPv4,
		ArpOp_REQUEST,
		h.config.localHrdAddr,
		h.config.localPrtAddr,
		[]byte{0, 0, 0, 0, 0, 0},
		dstPrtAddr,
	)

	if err := h.sender.ArpSend(nil, ap); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Resolve(targetPrtAddr []byte) ([]byte, error) {
	if addr, ok := globalArpTable[[4]byte(targetPrtAddr)]; ok {
		return addr[:], nil
	}

	if err := h.Request(targetPrtAddr); err != nil {
		return nil, err
	}

	addr := globalArpTable[[4]byte(targetPrtAddr)]

	return addr[:], nil
}
