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
	return &UnknownPayload{
		b: b,
	}
}

func (b *UnknownPayload) Bytes() []byte {
	return b.b
}

func (b *UnknownPayload) Inspect() {
	fmt.Println("Unknown Payload:")
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
