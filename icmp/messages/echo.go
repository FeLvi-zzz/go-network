package messages

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/icmp/types"
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type ICMPEchoPayload struct {
	Identifier     uint16
	SequenceNumber uint16
	Data           payload.Payload
}

func NewEchoReplyMessage(id uint16, seq uint16, d []byte) *ICMPMessage {
	m := &ICMPMessage{
		Type: types.ICMPType_EchoReply,
		Code: 0,
		Payload: &ICMPEchoPayload{
			Identifier:     id,
			SequenceNumber: seq,
			Data:           payload.NewDataPayload(d),
		},
	}

	m.Checksum = util.CalcCheckSum(m.Bytes())

	return m
}

func EchoPayloadFromBytes(b []byte) *ICMPEchoPayload {
	return &ICMPEchoPayload{
		Identifier:     util.ToUint16(b[0:2]),
		SequenceNumber: util.ToUint16(b[2:4]),
		Data:           payload.NewDataPayload(b[4:]),
	}
}

func (p *ICMPEchoPayload) Bytes() []byte {
	b := append(
		make([]byte, 0, 4),
		byte(p.Identifier>>8),
		byte(p.Identifier),
		byte(p.SequenceNumber>>8),
		byte(p.SequenceNumber),
	)
	b = append(b, p.Data.Bytes()...)

	return b
}

func (p *ICMPEchoPayload) Inspect() {
	fmt.Printf("  Identifier: %d\n", p.Identifier)
	fmt.Printf("  SequenceNumber: %d\n", p.SequenceNumber)
	p.Data.Inspect()
}
