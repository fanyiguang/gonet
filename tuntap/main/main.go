package main

import (
	"fmt"
	"net"
	"time"

	"github.com/gamexg/gonet/tuntap"
)

func main() {

	f := tuntap.NewTunTapF()

	// f.Remove()

	lAddr := net.IPv4(192, 168, 123, 11)
	network := net.IPv4(190, 168, 123, 0)
	mask := net.IPv4(255, 255, 255, 0)

	tap, err := f.NewTun(lAddr, network, mask)
	if err != nil {
		panic(err)
	}
	defer tap.Close()

	buf := make([]byte, 1024)
	for i := 0; i < 10; i++ {
		_, err := tap.Write(buf)
		if err != nil {
			panic(err)
		}
	}
	/*
		for i := 0; i < 10; i++ {
			b, err := tap.Read()
			if err != nil {
				panic(err)
			}
			fmt.Println(b)
		}*/

	fmt.Println("sleep")
	time.Sleep(30 * time.Second)

}
