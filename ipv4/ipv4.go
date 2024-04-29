package ipv4

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/ipv4/types"
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type IPv4Packet struct {
	Version        uint8 // 4bit
	IHL            uint8 // 4bit
	ToS            uint8
	TotalLength    uint16
	Identification uint16
	Flags          types.Flags // 3bit
	FlagmentOffset uint16      // 13bit
	TTL            uint8
	Protocol       types.Protocol // uint8
	HeaderChecksum uint16
	SrcAddr        types.Address
	DstAddr        types.Address
	Options        []byte

	Payload payload.Payload
}

func FromBytes(b []byte) (*IPv4Packet, []byte, error) {
	if len(b) < 20 {
		return nil, nil, fmt.Errorf("this is not IPv4 packet")
	}

	pkt := &IPv4Packet{
		Version:        (b[0] & 0xf0) >> 4,
		IHL:            b[0] & 0x0f,
		ToS:            b[1],
		TotalLength:    util.ToUint16(b[2:4]),
		Identification: util.ToUint16(b[4:6]),
		Flags:          types.Flags((b[6] & 0xe0) >> 5),
		FlagmentOffset: ((uint16(b[6]) & 0x1f) << 8) | uint16(b[7]),
		TTL:            b[8],
		Protocol:       types.Protocol(b[9]),
		HeaderChecksum: util.ToUint16(b[10:12]),
		SrcAddr:        types.Address(b[12:16]),
		DstAddr:        types.Address(b[16:20]),
	}

	dataStart := pkt.IHL * 4

	pkt.Options = b[20:dataStart]

	packetLength := int(pkt.TotalLength)
	if packetLength > len(b) {
		fmt.Printf("[DEBUG] packet length is short, TotalLength = %d, packetLength = %d\n", pkt.TotalLength, len(b))
		packetLength = len(b)
	}

	return pkt, b[dataStart:packetLength], nil
}

func (p *IPv4Packet) Inspect() {
	fmt.Println("IPv4 header:")
	fmt.Printf("  Version: %d\n", p.Version)
	fmt.Printf("  IHL: %d * 32bit\n", p.IHL)
	fmt.Printf("  ToS: %8b\n", p.ToS)
	fmt.Printf("  TotalLength: %d octets\n", p.TotalLength)
	fmt.Printf("  Identification: %d\n", p.Identification)
	fmt.Printf("  Flags: %s\n", p.Flags.ToString())
	fmt.Printf("  FlagmentOffset: %d\n", p.FlagmentOffset)
	fmt.Printf("  TTL: %d\n", p.TTL)
	fmt.Printf("  Protocol: %s\n", p.Protocol.ToString())
	fmt.Printf("  HeaderChecksum: %04x (isValid? => %t)\n", p.HeaderChecksum, p.IsValid())
	fmt.Printf("  SrcAddr: %s\n", p.SrcAddr.ToString())
	fmt.Printf("  DstAddr: %s\n", p.DstAddr.ToString())
	p.Payload.Inspect()
}

func (p *IPv4Packet) Bytes() []byte {
	return append(p.HeaderBytes(), p.Payload.Bytes()...)
}

func (p *IPv4Packet) HeaderBytes() []byte {
	b := make([]byte, 0, 20) // header min length

	b = append(
		b,
		p.Version<<4|p.IHL,
		p.ToS,
		byte(p.TotalLength>>8),
		byte(p.TotalLength),
		byte(p.Identification>>8),
		byte(p.Identification),
		uint8(p.Flags)<<5|uint8(p.FlagmentOffset>>8),
		byte(p.FlagmentOffset),
		p.TTL,
		byte(p.Protocol),
		byte(p.HeaderChecksum>>8),
		byte(p.HeaderChecksum),
	)
	b = append(b, p.SrcAddr[:]...)
	b = append(b, p.DstAddr[:]...)
	b = append(b, p.Options...)
	for i := 0; i < len(p.Options)%4; i++ {
		b = append(b, 0)
	}

	return b
}

func (p *IPv4Packet) IsValid() bool {
	return p.CalcChecksum() == 0
}

func (p *IPv4Packet) CalcChecksum() uint16 {
	return util.CalcCheckSum(p.HeaderBytes())
}
