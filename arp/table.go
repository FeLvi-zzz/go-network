package arp

type arpTable map[[4]byte][6]byte

var globalArpTable arpTable

func init() {
	globalArpTable = make(arpTable)
}
