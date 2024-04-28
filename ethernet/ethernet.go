package ethernet

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/ethernet/types"
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type EthernetFrame struct {
	DstMacAddr types.Address
	SrcMacAddr types.Address
	Ethertype  types.EtherType
	Payload    payload.Payload
}

func NewEthernetFrame(dst []byte, src []byte, etype types.EtherType, payload payload.Payload) *EthernetFrame {
	return &EthernetFrame{
		DstMacAddr: types.Address(dst),
		SrcMacAddr: types.Address(src),
		Ethertype:  etype,
		Payload:    payload,
	}
}

func (e *EthernetFrame) Bytes() []byte {
	b := make([]byte, 0, 1500) // MTU 1500

	b = append(b, e.DstMacAddr[:]...)
	b = append(b, e.SrcMacAddr[:]...)
	b = append(b, byte(e.Ethertype>>8), byte(e.Ethertype))
	b = append(b, e.Payload.Bytes()...)

	return b
}

func FromBytes(b []byte) (*EthernetFrame, []byte, error) {
	if len(b) < 14 {
		return nil, nil, fmt.Errorf("this is not ethernet packet")
	}

	return &EthernetFrame{
		DstMacAddr: types.Address(b[0:6]),
		SrcMacAddr: types.Address(b[6:12]),
		Ethertype:  types.EtherType(util.ToUint16(b[12:14])),
	}, b[14:], nil
}

func (e *EthernetFrame) Inspect() {
	fmt.Println("Ethernet Frame:")
	fmt.Printf("  dest: %s\n", e.DstMacAddr.ToString())
	fmt.Printf("  src: %s\n", e.SrcMacAddr.ToString())
	fmt.Printf("  EtherType: %s\n", e.Ethertype.ToString())
	e.Payload.Inspect()
}
