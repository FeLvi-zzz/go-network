package types

import "fmt"

type Address [4]byte

func (a Address) ToString() string {
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}

func (a Address) CalcSubnet(subnetMask uint8) Address {
	s := a
	sm := ^uint32(0) << (32 - subnetMask)
	for i := range s {
		s[i] &= uint8(sm >> (32 - 8*(i+1)))
	}
	return s
}
