package arp

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/ethernet"
	"github.com/FeLvi-zzz/go-network/util"
)

type ArpOp uint16

const (
	ArpOp_REQUEST  ArpOp = 1
	ArpOp_RESPONSE ArpOp = 2
)

func (a ArpOp) ToString() string {
	switch a {
	case ArpOp_REQUEST:
		return "REQUEST"
	case ArpOp_RESPONSE:
		return "RESPONSE"
	default:
		return "unknown"
	}
}

type ArpPayload struct {
	Hrd uint16
	Pro ethernet.EtherType // uint16
	Hln uint8
	Pln uint8
	Op  ArpOp
	Sha []byte
	Spa []byte
	Tha []byte
	Tpa []byte
}

func NewPayload(
	protocol ethernet.EtherType,
	op ArpOp,
	localHardAddr []byte,
	localProtoAddr []byte,
	remoteHardAddr []byte,
	remoteProtoAddr []byte,
) *ArpPayload {
	switch protocol {
	case ethernet.EtherType_IPv4:
		return &ArpPayload{
			Hrd: 1, // Ethernet = 1
			Pro: ethernet.EtherType_IPv4,
			Hln: 6, // Ethernet address length
			Pln: 4, // IPv4 address length
			Op:  op,
			Sha: localHardAddr,
			Spa: localProtoAddr,
			Tha: remoteHardAddr,
			Tpa: remoteProtoAddr,
		}
	default:
		return &ArpPayload{}
	}
}

func (p *ArpPayload) ToBytes() []byte {
	b := make([]byte, 0, 8+p.Hln+p.Pln+p.Hln+p.Pln)

	b = append(b, []byte{
		uint8(p.Hrd >> 8),
		uint8(p.Hrd),
		uint8(p.Pro >> 8),
		uint8(p.Pro),
		p.Hln,
		p.Pln,
		uint8(p.Op >> 8),
		uint8(p.Op),
	}...)
	b = append(b, p.Sha...)
	b = append(b, p.Spa...)
	b = append(b, p.Tha...)
	b = append(b, p.Tpa...)

	return b
}

func FromBytes(b []byte) (*ArpPayload, error) {
	if len(b) < 8 {
		return nil, fmt.Errorf("this arp payload is broken")
	}

	hln := uint8(b[4])
	pln := uint8(b[5])

	if len(b) < int(8+2*(hln+pln)) {
		return nil, fmt.Errorf("this arp payload is broken")
	}

	return &ArpPayload{
		Hrd: util.ToUint16(b[0:2]),
		Pro: ethernet.EtherType(util.ToUint16(b[2:4])),
		Hln: hln,
		Pln: pln,
		Op:  ArpOp(util.ToUint16(b[6:8])),
		Sha: b[8 : 8+hln],
		Spa: b[8+hln : 8+hln+pln],
		Tha: b[8+hln+pln : 8+hln+pln+hln],
		Tpa: b[8+hln+pln+hln : 8+hln+pln+hln+pln],
	}, nil
}

func (a *ArpPayload) Inspect() {
	fmt.Println("ARP payload:")
	fmt.Printf("  HardwareType: %x\n", a.Hrd)
	fmt.Printf("  Protocol: %s\n", a.Pro.ToString())
	fmt.Printf("  Op: %s\n", a.Op.ToString())
	fmt.Printf("  SrcHardwareAddress: %s\n", ethernet.HardwareAddressToString(a.Sha))
	fmt.Printf("  SrcProtocolAddress: %s\n", ProtocolAddressToString(a.Spa, a.Pro))
	fmt.Printf("  DstHardwareAddress: %s\n", ethernet.HardwareAddressToString(a.Tha))
	fmt.Printf("  DstProtocolAddress: %s\n", ProtocolAddressToString(a.Tpa, a.Pro))
}

func ProtocolAddressToString(b []byte, protocol ethernet.EtherType) string {
	switch protocol {
	case ethernet.EtherType_IPv4:
		if len(b) != 4 {
			panic("this is not ipv4 address")
		}
		return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
	default:
		panic(fmt.Sprintf("unknown protocol: %s", protocol.ToString()))
	}
}
