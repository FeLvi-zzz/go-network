package types

import "fmt"

type Protocol uint8

const (
	Protocol_ICMP Protocol = 1
	Protocol_TCP  Protocol = 6
	Protocol_UDP  Protocol = 17
)

func (p Protocol) ToString() string {
	switch p {
	case Protocol_ICMP:
		return "ICMP"
	case Protocol_TCP:
		return "TCP"
	case Protocol_UDP:
		return "UDP"
	}

	return fmt.Sprintf("Unknown(%d)", p)
}
