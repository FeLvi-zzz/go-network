package main

import (
	"syscall"

	"github.com/FeLvi-zzz/go-network/arp"
	"github.com/FeLvi-zzz/go-network/ethernet"
	"github.com/FeLvi-zzz/go-network/icmp"
	"github.com/FeLvi-zzz/go-network/ipv4"
	ipv4types "github.com/FeLvi-zzz/go-network/ipv4/types"
	"github.com/FeLvi-zzz/go-network/kernel"
	"github.com/FeLvi-zzz/go-network/tcp"
	tcpsample "github.com/FeLvi-zzz/go-network/tcp/sample"
	"github.com/FeLvi-zzz/go-network/udp"
	udpsample "github.com/FeLvi-zzz/go-network/udp/sample"
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

	srcHrdAddr := []byte{0x00, 0x15, 0x5d, 0xc5, 0x77, 0xf3}
	broadcastHrdAddr := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	srcPrtAddr := []byte{172, 20, 159, 90}
	dstPrtAddr := []byte{172, 20, 144, 1}
	ifIndex := 2

	kernelConfig := kernel.NewConfig(fd, ifIndex)
	ethConfig := ethernet.NewConfig(srcHrdAddr, broadcastHrdAddr)
	arpConfig := arp.NewConfig(srcHrdAddr, srcPrtAddr)
	ipv4Config := ipv4.NewConfig(srcPrtAddr, ipv4.NewRouteTable(
		[]ipv4.RouteTableRecord{
			{Subnet: ipv4types.Address(dstPrtAddr), SubnetMask: 20, TargetIp: ipv4types.Address{}},
			{Subnet: ipv4types.Address{}, SubnetMask: 0, TargetIp: ipv4types.Address(dstPrtAddr)},
		},
	))
	icmpConfig := icmp.NewConfig()
	udpConfig := udp.NewConfig()
	tcpConfig := tcp.NewConfig()

	kernelSender := kernel.NewSender(kernelConfig)
	ethSender := ethernet.NewSender(ethConfig, kernelSender)

	arpHandler := arp.NewHandler(arpConfig, ethSender)

	ipv4Sender := ipv4.NewSender(ipv4Config, ethSender, arpHandler)
	icmpConsumer := icmp.NewConsumer(icmpConfig, ipv4Sender)
	arpConsumer := arp.NewConsumer(arpConfig, ethSender)
	udpConsumer := udp.NewConsumer(udpConfig, ipv4Sender)
	tcpConsumer := tcp.NewConsumer(tcpConfig, ipv4Sender)

	ipv4Consumer := ipv4.NewConsumer(ipv4Config, ethSender, icmpConsumer, udpConsumer, tcpConsumer)
	ethConsumer := ethernet.NewConsumer(ethConfig, arpConsumer, ipv4Consumer, kernelSender)

	kernelHandler := kernel.NewHandler(kernelConfig, ethConsumer)

	errch := make(chan error)
	go func() {
		errch <- kernelHandler.Handle()
	}()
	go func() {
		_, err := arpHandler.Resolve(dstPrtAddr)
		errch <- err
	}()
	go func() {
		errch <- udpsample.Serve(ipv4Sender, srcPrtAddr, 3000)
	}()
	go func() {
		errch <- udpsample.RequestHoge(ipv4Sender, dstPrtAddr, 4000, srcPrtAddr, 4000)
	}()
	go func() {
		errch <- tcpsample.Serve(ipv4Sender, srcPrtAddr, 5000)
	}()

	for {
		if err := <-errch; err != nil {
			return err
		}
	}
}

func htons(i uint16) uint16 {
	return (i&0xff)<<8 | (i >> 8)
}
