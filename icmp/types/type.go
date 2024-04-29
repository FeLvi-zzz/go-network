package types

import "fmt"

type ICMPType uint8

const (
	ICMPType_EchoReply ICMPType = 0
	ICMPType_Echo      ICMPType = 8
)

func (i ICMPType) ToString() string {
	switch i {
	case ICMPType_EchoReply:
		return "Echo Reply"
	case ICMPType_Echo:
		return "Echo"
	}

	return fmt.Sprintf("Unknown(%d)", i)
}
