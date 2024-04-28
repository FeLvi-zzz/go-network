package ipv4

import (
	"fmt"

	"github.com/FeLvi-zzz/go-network/payload"
)

type Consumer struct {
	config          *Config
	fragmentBuilder fragmentBuilder
	// sender sender
}

type fragmentBuilder map[uint16]map[int][]byte

func (f fragmentBuilder) Add(id uint16, offset int, b []byte) {
	if f[id] == nil {
		f[id] = make(map[int][]byte)
	}
	f[id][offset] = b
}

func (f fragmentBuilder) Build(id uint16) ([]byte, error) {
	defer delete(f, id)

	l := 0
	for _, b := range f[id] {
		l += len(b)
	}

	res := make([]byte, 0, l)
	offset := 0
	for {
		b := f[id][offset]
		delete(f[id], offset)

		res = append(res, b...)

		if b == nil {
			if len(f[id]) > 0 {
				return nil, fmt.Errorf("broken fragments")
			}
			return res, nil
		}

		offset = len(res)
	}
}

func NewConsumer(config *Config) *Consumer {
	return &Consumer{
		config:          config,
		fragmentBuilder: make(fragmentBuilder),
	}
}

func (c *Consumer) Consume(b []byte) (payload.Payload, error) {
	v4p, rb, err := FromBytes(b)
	if err != nil {
		return payload.NewUnknownPayload(b), err
	}

	c.fragmentBuilder.Add(v4p.Identification, int(v4p.FlagmentOffset)*8, rb)

	if v4p.Flags.IsMF() {
		v4p.Payload = payload.NewFragmentPayload(len(rb))
	} else {
		nrb, err := c.fragmentBuilder.Build(v4p.Identification)
		if err != nil {
			return payload.NewFragmentPayload(len(rb)), err
		}
		// TODO: consume higher protocol
		// up, err := c.udpConsumer.Consume(c.fragmentBuffer[v4p.Identification])
		// if err != nil {
		// 	return payload.NewUnknownPayload(c.fragmentBuffer[v4p.Identification]), err
		// }
		// v4p.Payload = up

		v4p.Payload = payload.NewUnknownPayload(nrb)
	}

	return v4p, nil
}
