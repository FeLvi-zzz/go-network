package util

func ToUint16(b []byte) uint16 {
	if len(b) != 2 {
		panic("must be 2 bytes")
	}
	return (uint16(b[0]) << 8) | uint16(b[1])
}
