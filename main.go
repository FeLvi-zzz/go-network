package main

import (
	"sync"
	"syscall"

	"github.com/FeLvi-zzz/go-network/arp"
)

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}

func _main() error {
	// create socket
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		return err
	}

	srcHrdAddr := []byte{0x00, 0x15, 0x5d, 0x55, 0xaa, 0xfd}
	broadcastHrdAddr := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	srcPrtAddr := []byte{172, 20, 159, 90}
	dstPrtAddr := []byte{172, 20, 144, 1}
	ifIndex := 2

	mu := new(sync.Mutex)

	arpsvc := arp.NewService(mu, fd, srcHrdAddr, broadcastHrdAddr, srcPrtAddr, dstPrtAddr, ifIndex)

	errch := make(chan error)
	go func() {
		errch <- arpsvc.Start()
	}()

	return <-errch
}

func htons(i uint16) uint16 {
	return (i&0xff)<<8 | (i >> 8)
}
