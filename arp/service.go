package arp

import (
	"fmt"
	"time"

	ethtypes "github.com/FeLvi-zzz/go-network/ethernet/types"
)

type Conn struct {
	res chan []byte
}

var conn *Conn

func init() {
	conn = &Conn{
		res: make(chan []byte, 10),
	}
}

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

func (h *Handler) request(dstPrtAddr []byte) ([]byte, error) {
	ap := NewPayload(
		ethtypes.EtherType_IPv4,
		ArpOp_REQUEST,
		h.config.localHrdAddr,
		h.config.localPrtAddr,
		[]byte{0, 0, 0, 0, 0, 0},
		dstPrtAddr,
	)

	if err := h.sender.ArpSend(nil, ap); err != nil {
		return nil, err
	}

	timer := time.NewTimer(1 * time.Second)

	select {
	case <-timer.C:
		return nil, fmt.Errorf("arp timeout")
	case addr := <-conn.res:
		return addr, nil
	}
}

func (h *Handler) Resolve(targetPrtAddr []byte) ([]byte, error) {
	if addr, ok := globalArpTable[[4]byte(targetPrtAddr)]; ok {
		return addr[:], nil
	}

	addr, err := h.request(targetPrtAddr)
	if err != nil {
		return nil, err
	}

	return addr[:], nil
}
