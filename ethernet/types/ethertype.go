package types

import (
	"fmt"
	"syscall"
)

type EtherType uint16

const (
	EtherType_IPv4 = syscall.ETH_P_IP
	EtherType_ARP  = syscall.ETH_P_ARP
)

func (t EtherType) ToString() string {
	switch t {
	case EtherType_ARP:
		return "ARP"
	case EtherType_IPv4:
		return "IPv4"
	default:
		return fmt.Sprintf("Unknown(%x)", t)
	}
}
