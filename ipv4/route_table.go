package ipv4

import (
	"github.com/FeLvi-zzz/go-network/ipv4/types"
)

type RouteTable []RouteTableRecord
type RouteTableRecord struct {
	Subnet     types.Address
	SubnetMask uint8
	TargetIp   types.Address
}

func NewRouteTable(records []RouteTableRecord) RouteTable {
	return records
}

func (t RouteTable) Resolve(tip types.Address) types.Address {
	for _, r := range t {
		ips := r.Subnet.CalcSubnet(r.SubnetMask)
		tips := tip.CalcSubnet(r.SubnetMask)
		if ips == tips {
			if r.TargetIp == types.Address([4]byte{}) {
				return tip
			}
			return r.TargetIp
		}
	}

	return types.Address{}
}
