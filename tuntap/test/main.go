package main

import (
	"github.com/gamexg/gonet/tuntap"
	"net"
	"time"
)

// 尝试手工实现 tcp 协议，这样测试一下性能。


func main() {
	f := tuntap.NewTunTapF()
	lAddr := net.IPv4(192, 168, 150, 11)
	network := net.IPv4(190, 168, 150, 0)
	mask := net.IPv4(255, 255, 255, 0)

	tap, err := f.NewTun(lAddr, network, mask)
	if err != nil {
		panic(err)
	}
	defer tap.Close()


	time.Sleep()

}


