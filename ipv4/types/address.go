package types

import "fmt"

type Address [4]byte

func (a Address) ToString() string {
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}
