package types

import "fmt"

type Address [6]byte

func (a Address) ToString() string {
	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		a[0],
		a[1],
		a[2],
		a[3],
		a[4],
		a[5],
	)
}
