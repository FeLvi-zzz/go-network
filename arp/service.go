package arp

import (
	"bytes"
	"fmt"
	"sync"
	"syscall"

	"github.com/FeLvi-zzz/go-network/ethernet"
)

type service struct {
	mu               *sync.Mutex
	Fd               int
	LocalHrdAddr     []byte
	BroadcastHrdAddr []byte
	LocalPrtAddr     []byte
	RemotePrtAddr    []byte
	IfIndex          int
}

func NewService(
	mu *sync.Mutex,
	fd int,
	localHrdAddr []byte,
	broadcastHrdAddr []byte,
	localPrtAddr []byte,
	remotePrtAddr []byte,
	ifIndex int,
) *service {
	return &service{
		mu:               mu,
		Fd:               fd,
		LocalHrdAddr:     localHrdAddr,
		BroadcastHrdAddr: broadcastHrdAddr,
		LocalPrtAddr:     localPrtAddr,
		RemotePrtAddr:    remotePrtAddr,
		IfIndex:          ifIndex,
	}
}

func (s *service) Start() error {
	if err := s.send(nil); err != nil {
		return err
	}

	for {
		ap, err := s.recv()
		if err != nil {
			return err
		}
		if ap == nil {
			continue
		}

		if ap.Op == ArpOp_REQUEST && bytes.Equal(ap.Tpa, s.LocalPrtAddr) {
			if err := s.send(ap.Sha); err != nil {
				return err
			}
		}
	}
}

func (s *service) send(remoteHrdAddr []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	packet := make([]byte, 0, 1500)

	var (
		ethDst    []byte
		arpRemote []byte
		op        ArpOp
	)

	if remoteHrdAddr == nil {
		ethDst = s.BroadcastHrdAddr
		arpRemote = []byte{0, 0, 0, 0, 0, 0}
		op = ArpOp_REQUEST
	} else {
		ethDst = remoteHrdAddr
		arpRemote = remoteHrdAddr
		op = ArpOp_RESPONSE
	}

	eh := ethernet.NewHeader(
		ethDst,
		s.LocalHrdAddr,
		ethernet.EtherType_ARP,
	)
	ap := NewPayload(
		ethernet.EtherType_IPv4,
		op,
		s.LocalHrdAddr,
		s.LocalPrtAddr,
		arpRemote,
		s.RemotePrtAddr,
	)

	packet = append(packet, eh.ToBytes()...)
	packet = append(packet, ap.ToBytes()...)

	fmt.Printf("\n-- send packet --\n")
	eh.Inspect()
	ap.Inspect()

	if err := syscall.Sendto(s.Fd, packet, syscall.MSG_CONFIRM, &syscall.SockaddrLinklayer{
		Ifindex: s.IfIndex,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) recv() (*ArpPayload, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b := make([]byte, 80)
	_, _, err := syscall.Recvfrom(s.Fd, b, 0)
	if err != nil {
		return nil, err
	}

	ef, rb, err := ethernet.FromBytes(b)
	if err != nil {
		return nil, err
	}

	if ef.Ethertype == ethernet.EtherType_ARP {
		ap, err := FromBytes(rb)
		if err != nil {
			return nil, err
		}
		fmt.Printf("\n-- recv packet --\n")
		ef.Inspect()
		ap.Inspect()
		fmt.Println("")

		return ap, nil
	}

	return nil, nil
}
