package util

func ToUint16(b []byte) uint16 {
	if len(b) != 2 {
		panic("must be 2 bytes")
	}
	return (uint16(b[0]) << 8) | uint16(b[1])
}

func ToUint32(b []byte) uint32 {
	if len(b) != 4 {
		panic("must be 4 bytes")
	}
	return (uint32(ToUint16(b[0:2])) << 16) | uint32(ToUint16(b[2:4]))
}
