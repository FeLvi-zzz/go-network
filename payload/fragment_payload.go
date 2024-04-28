package payload

import "fmt"

type FragmentPayload struct {
	len int
}

func NewFragmentPayload(len int) *FragmentPayload {
	return &FragmentPayload{
		len: len,
	}
}

func (b *FragmentPayload) Bytes() []byte {
	return []byte{}
}

func (b *FragmentPayload) Inspect() {
	fmt.Println("Fragment Payload:")
	fmt.Printf("  length: %d bytes\n", b.len)
}
