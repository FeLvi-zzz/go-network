package types

import "strings"

type Flags uint8

const (
	Flags_CWR Flags = 1 << (7 - iota)
	Flags_ECE Flags = 1 << (7 - iota)
	Flags_URG Flags = 1 << (7 - iota)
	Flags_ACK Flags = 1 << (7 - iota)
	Flags_PSH Flags = 1 << (7 - iota)
	Flags_RST Flags = 1 << (7 - iota)
	Flags_SYN Flags = 1 << (7 - iota)
	Flags_FIN Flags = 1 << (7 - iota)
)

func (f Flags) ToString() string {
	ss := make([]string, 0, 8)

	if f&Flags_CWR != 0 {
		ss = append(ss, "CWR")
	}
	if f&Flags_ECE != 0 {
		ss = append(ss, "ECE")
	}
	if f&Flags_URG != 0 {
		ss = append(ss, "URG")
	}
	if f&Flags_ACK != 0 {
		ss = append(ss, "ACK")
	}
	if f&Flags_PSH != 0 {
		ss = append(ss, "PSH")
	}
	if f&Flags_RST != 0 {
		ss = append(ss, "RST")
	}
	if f&Flags_SYN != 0 {
		ss = append(ss, "SYN")
	}
	if f&Flags_FIN != 0 {
		ss = append(ss, "FIN")
	}

	return strings.Join(ss, "|")
}
