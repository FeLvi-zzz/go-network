package tcp

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/payload"
	"github.com/FeLvi-zzz/go-network/tcp/types"
	"github.com/FeLvi-zzz/go-network/util"
)

type Segment struct {
	SrcPort    uint16
	DstPort    uint16
	SeqNum     uint32
	AckNum     uint32
	DataOffset uint8       // 4bit
	Flags      types.Flags // 8bit
	Window     uint16
	Checksum   uint16
	UrgPtr     uint16
	Options    []byte
	Data       payload.Payload

	PseudoHeader []byte
}

func FromBytes(b []byte, ph []byte) (*Segment, error) {
	if len(b) < 20 {
		return nil, fmt.Errorf("tcp header is broken")
	}

	dataOffset := b[12] >> 4

	return &Segment{
		SrcPort:    util.ToUint16(b[0:2]),
		DstPort:    util.ToUint16(b[2:4]),
		SeqNum:     util.ToUint32(b[4:8]),
		AckNum:     util.ToUint32(b[8:12]),
		DataOffset: dataOffset,
		Flags:      types.Flags(b[13]),
		Window:     util.ToUint16(b[14:16]),
		Checksum:   util.ToUint16(b[16:18]),
		UrgPtr:     util.ToUint16(b[18:20]),
		Options:    b[20 : dataOffset*4],
		Data:       payload.NewDataPayload(b[dataOffset*4:]),

		PseudoHeader: ph,
	}, nil
}

func (s *Segment) Bytes() []byte {
	b := make([]byte, 0, 20)
	b = append(
		b,
		byte(s.SrcPort>>8),
		byte(s.SrcPort),
		byte(s.DstPort>>8),
		byte(s.DstPort),
		byte(s.SeqNum>>24),
		byte(s.SeqNum>>16),
		byte(s.SeqNum>>8),
		byte(s.SeqNum),
		byte(s.AckNum>>24),
		byte(s.AckNum>>16),
		byte(s.AckNum>>8),
		byte(s.AckNum),
		byte(s.DataOffset<<4), // includes reserved field
		byte(s.Flags),
		byte(s.Window>>8),
		byte(s.Window),
		byte(s.Checksum>>8),
		byte(s.Checksum),
		byte(s.UrgPtr>>8),
		byte(s.UrgPtr),
	)
	b = append(b, s.Options...)
	b = append(b, s.Data.Bytes()...)

	return b
}

func (s *Segment) Inspect() {
	fmt.Println("TCP Segment:")
	fmt.Printf("  SrcPort: %d\n", s.SrcPort)
	fmt.Printf("  DstPort: %d\n", s.DstPort)
	fmt.Printf("  SeqNum: %d\n", s.SeqNum)
	fmt.Printf("  AckNum: %d\n", s.AckNum)
	fmt.Printf("  DataOffset: %d * 32bit\n", s.DataOffset)
	fmt.Printf("  Flags: %s\n", s.Flags.ToString())
	fmt.Printf("  Window: %d\n", s.Window)
	fmt.Printf("  Checksum: %04x (isValid? => %t)\n", s.Checksum, s.IsValid())
	fmt.Printf("  UrgPtr: %d\n", s.UrgPtr)
	fmt.Printf("  Options: %#v\n", s.Options)
	s.Data.Inspect()
}

func (s *Segment) CalcCheckSum() uint16 {
	return util.CalcCheckSum(append(s.PseudoHeader, s.Bytes()...))
}

func (s *Segment) IsValid() bool {
	return s.CalcCheckSum() == 0
}

func NewSegment(srcPort uint16, dstPort uint16, seqNum uint32, ackNum uint32, flags types.Flags, data payload.Payload) *Segment {
	return &Segment{
		SrcPort:    srcPort,
		DstPort:    dstPort,
		SeqNum:     seqNum,
		AckNum:     ackNum,
		DataOffset: 5,
		Flags:      flags,
		Window:     16384,
		// Checksum:   0 // calc later,
		// UrgPtr:     0,
		Data: data,
	}
}
