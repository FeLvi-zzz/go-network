package udp

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/util"
)

type Datagram struct {
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
	Data     payload.Payload

	PseudoHeader [12]byte
}

func FromBytes(b []byte, ph []byte) (*Datagram, error) {
	if len(b) < 8 {
		return nil, fmt.Errorf("udp header is broken")
	}

	return &Datagram{
		SrcPort:      util.ToUint16(b[0:2]),
		DstPort:      util.ToUint16(b[2:4]),
		Length:       util.ToUint16(b[4:6]),
		Checksum:     util.ToUint16(b[6:8]),
		Data:         payload.NewDataPayload(b[8:]),
		PseudoHeader: [12]byte(ph),
	}, nil
}

func (d *Datagram) Bytes() []byte {
	b := make([]byte, 0, 8)
	b = append(
		b,
		byte(d.SrcPort>>8),
		byte(d.SrcPort),
		byte(d.DstPort>>8),
		byte(d.DstPort),
		byte(d.Length>>8),
		byte(d.Length),
		byte(d.Checksum>>8),
		byte(d.Checksum),
	)
	b = append(b, d.Data.Bytes()...)

	return b
}

func (d *Datagram) Inspect() {
	fmt.Println("UDP Datagram:")
	fmt.Printf("  SrcPort: %d\n", d.SrcPort)
	fmt.Printf("  DstPort: %d\n", d.DstPort)
	fmt.Printf("  Length: %d\n", d.Length)
	fmt.Printf("  Checksum: %04x (isValid? => %t)\n", d.Checksum, d.IsValid())
	d.Data.Inspect()
}

func (d *Datagram) IsValid() bool {
	return util.CalcCheckSum(append(d.PseudoHeader[:], d.Bytes()...)) == 0
}
