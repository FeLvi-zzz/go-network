package ethernet

import (
	"fmt"
	"syscall"

	"github.com/FeLvi-zzz/go-network/util"
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

type EthernetFrame struct {
	DstMacAddr [6]byte
	SrcMacAddr [6]byte
	Ethertype  EtherType
}

func NewHeader(dst []byte, src []byte, etype uint16) *EthernetFrame {
	return &EthernetFrame{
		DstMacAddr: [6]byte(dst),
		SrcMacAddr: [6]byte(src),
		Ethertype:  EtherType(etype),
	}
}

func (e *EthernetFrame) ToBytes() []byte {
	b := make([]byte, 0, 14)
	b = append(b, e.DstMacAddr[:]...)
	b = append(b, e.SrcMacAddr[:]...)
	b = append(b, byte(e.Ethertype>>8), byte(e.Ethertype))

	return b
}

func FromBytes(b []byte) (*EthernetFrame, []byte, error) {
	if len(b) < 14 {
		return nil, nil, fmt.Errorf("this is not ethernet packet")
	}

	return &EthernetFrame{
		DstMacAddr: [6]byte(b[0:6]),
		SrcMacAddr: [6]byte(b[6:12]),
		Ethertype:  EtherType(util.ToUint16(b[12:14])),
	}, b[14:], nil
}

func (e *EthernetFrame) Inspect() {
	fmt.Println("Ethernet Frame:")
	fmt.Printf("  dest: %s\n", HardwareAddressToString(e.DstMacAddr[:]))
	fmt.Printf("  src: %s\n", HardwareAddressToString(e.SrcMacAddr[:]))
	fmt.Printf("  Ethertype: %s\n", e.Ethertype.ToString())
}

func HardwareAddressToString(b []byte) string {
	if len(b) != 6 {
		panic("this is not hardware address")
	}

	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		b[0],
		b[1],
		b[2],
		b[3],
		b[4],
		b[5],
	)
}
