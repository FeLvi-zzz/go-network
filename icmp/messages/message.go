package messages

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/icmp/types"
	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type ICMPMessage struct {
	Type     types.ICMPType // uint8
	Code     uint8
	Checksum uint16
	Payload  payload.Payload
}

func (m *ICMPMessage) Bytes() []byte {
	b := append(
		make([]byte, 0, 8),
		byte(m.Type),
		m.Code,
		byte(m.Checksum>>8),
		byte(m.Checksum),
	)
	b = append(b, m.Payload.Bytes()...)

	return b
}

func (m *ICMPMessage) Inspect() {
	fmt.Println("ICMP Message:")
	fmt.Printf("  Type: %s\n", m.Type.ToString())
	fmt.Printf("  Code: %d\n", m.Code)
	fmt.Printf("  Checksum: %04x (isValid? => %t)\n", m.Checksum, m.IsValid())
	m.Payload.Inspect()
}

func (m *ICMPMessage) IsValid() bool {
	return util.CalcCheckSum(m.Bytes()) == 0
}

func FromBytes(b []byte) (*ICMPMessage, error) {
	m := &ICMPMessage{
		Type:     types.ICMPType(b[0]),
		Code:     b[1],
		Checksum: util.ToUint16(b[2:4]),
	}

	switch m.Type {
	case types.ICMPType_Echo:
		m.Payload = EchoPayloadFromBytes(b[4:])
	case types.ICMPType_EchoReply:
		m.Payload = EchoPayloadFromBytes(b[4:])
	default:
		m.Payload = payload.NewUnknownPayload(b[4:])
	}

	return m, nil
}
