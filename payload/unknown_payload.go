package payload

import (
	"errors"
	"fmt"
)

var ErrUnknownPayload = errors.New("unknown payload")

type UnknownPayload struct {
	b []byte
}

func NewUnknownPayload(b []byte) *UnknownPayload {
	nb := make([]byte, len(b))
	copy(nb, b)

	return &UnknownPayload{
		b: nb,
	}
}

func (b *UnknownPayload) Bytes() []byte {
	return b.b
}

func (b *UnknownPayload) Inspect() {
	fmt.Println("Unknown Payload:")
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
