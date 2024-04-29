package payload

import (
	"fmt"
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
	for i := 0; i < len(b.b); i++ {
		if i%16 == 0 {
			fmt.Print("  ")
		}

		fmt.Printf("%02x", b.b[i])

		if i%16 == 15 {
			fmt.Printf("\n")
		} else {
			fmt.Print(" ")
		}
	}
	if len(b.b)%16 != 0 {
		fmt.Println("")
	}
}
