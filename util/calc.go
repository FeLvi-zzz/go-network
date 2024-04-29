package util

import "math"

func CalcCheckSum(b []byte) uint16 {
	if len(b)%2 == 1 {
		b = append(b, 0)
	}

	var cs uint16 = 0
	for i := 0; i < len(b); i += 2 {
		cs = AddOnesComplementUint16(cs, ToUint16(b[i:i+2]))
	}

	return ^cs
}

func AddOnesComplementUint16(a, b uint16) uint16 {
	if math.MaxUint16-a < b {
		return a + b + 1
	} else {
		return a + b
	}
}
