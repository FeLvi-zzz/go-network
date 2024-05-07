package payload

import (
	"fmt"
	"strings"
)

type DataPayload struct {
	b []byte
}

func NewDataPayload(b []byte) *DataPayload {
	nb := make([]byte, len(b))
	copy(nb, b)

	return &DataPayload{
		b: nb,
	}
}

func (b *DataPayload) Bytes() []byte {
	return b.b
}

func (b *DataPayload) Inspect() {
	fmt.Println("Data:")
	fmt.Printf("  length: %d bytes\n", len(b.b))

	if len(b.b) == 0 {
		return
	}

	hexstr := make([][]string, len(b.b)/16+1)
	asciistr := make([][]byte, len(b.b)/16+1)

	for i := 0; i < len(b.b); i++ {
		if i%16 == 0 {
			hexstr[i/16] = make([]string, 16)
			asciistr[i/16] = make([]byte, 16)
		}

		hexstr[i/16][i%16] = fmt.Sprintf("%02x", b.b[i])
		var as byte
		if 0x20 <= b.b[i] && b.b[i] <= 0x7e {
			as = b.b[i]
		} else {
			as = byte('.')
		}
		asciistr[i/16][i%16] = as
	}

	for i := 0; i < len(b.b)/16+1; i++ {
		for j := range hexstr[i] {
			if hexstr[i][j] == "" {
				hexstr[i][j] = "  "
			}
		}
		fmt.Printf("  %s | %s\n", strings.Join(hexstr[i], " "), asciistr[i])
	}
}
