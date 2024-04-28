package main

import (
	"syscall"

	"github.com/FeLvi-zzz/go-network/arp"
	"github.com/FeLvi-zzz/go-network/ethernet"
	"github.com/FeLvi-zzz/go-network/kernel"
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

	kernelConfig := kernel.NewConfig(fd, ifIndex)
	ethConfig := ethernet.NewConfig(srcHrdAddr, broadcastHrdAddr)
	arpConfig := arp.NewConfig(srcHrdAddr, srcPrtAddr)

	kernelSender := kernel.NewSender(kernelConfig)
	ethSender := ethernet.NewSender(ethConfig, kernelSender)

	arpConsumer := arp.NewConsumer(arpConfig, ethSender)
	ethConsumer := ethernet.NewConsumer(ethConfig, arpConsumer, kernelSender)

	kernelHandler := kernel.NewHandler(kernelConfig, ethConsumer)
	arpHandler := arp.NewHandler(arpConfig, ethSender)

	if err := arpHandler.Startup(dstPrtAddr); err != nil {
		return err
	}

	errch := make(chan error)
	go func() {
		errch <- kernelHandler.Handle()
	}()

	return <-errch
}

func htons(i uint16) uint16 {
	return (i&0xff)<<8 | (i >> 8)
}
