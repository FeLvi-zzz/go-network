package payload

type Payload interface {
	Bytes() []byte
	Inspect()
}
