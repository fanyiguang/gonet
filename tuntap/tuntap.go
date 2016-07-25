package tuntap

import (
	"net"

	"github.com/gamexg/gonet/forward"
)

type TunTapF interface {
	Remove() error
	NewTun(addr net.IP, network net.IP, mask net.IP) (Tun, error)
}

type Tun interface {
	forward.WritePacker
	forward.ReadPacker
	GetWriteChan() chan []byte
	GetReadChan() chan []byte
	Server() error
	Close() error
	GetName() string
}

type tunTapF struct {
}

func NewTunTapF() TunTapF {
	return &tunTapF{}
}
