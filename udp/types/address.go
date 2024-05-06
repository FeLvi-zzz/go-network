package types

import (
	"fmt"

	ipv4types "github.com/FeLvi-zzz/go-network/ipv4/types"
)

type Address struct {
	IP   []byte
	Port uint16
}

func (a Address) String() string {
	if len(a.IP) == 4 {
		return fmt.Sprintf("%s:%d", ipv4types.Address(a.IP).ToString(), a.Port)
	}

	panic(fmt.Errorf("unknown address type, '%s'", a.IP))
}
