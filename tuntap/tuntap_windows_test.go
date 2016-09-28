package tuntap

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"sync"

	"github.com/gamexg/gotool/mem"
	"github.com/google/gopacket/pcap"
)

const testDeviceName = "tap0901"

func TestNewAndRemove(t *testing.T) {
	/*
		f := &tunTapF{}
		err := f.newTunDevice(testDeviceName)
		if err != nil {
			t.Fatal(err)
		}

		id, err := getTuntapComponentId(testDeviceName)
		if err != nil {
			t.Fatal(err)
		}

		if id == "" {
			t.Fatal(`id==""`)
		}

		err = f.Remove()
		if err != nil {
			t.Fatal(err)
		}

		_, err = getTuntapComponentId(testDeviceName)
		if err == nil {
			t.Fatal("[err]remove")
		}
	*/
}

func TestAll(t *testing.T) {
	f := NewTunTapF()
	lAddr := net.IPv4(192, 168, 150, 11)
	network := net.IPv4(190, 168, 150, 0)
	mask := net.IPv4(255, 255, 255, 0)

	tap, err := f.NewTun(lAddr, network, mask)
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()
	//defer f.Remove()

	func() {
		buf := make([]byte, 1024)
		_, err := tap.Write(buf)
		if err != nil {
			t.Fatal(err)
		}

	}()

	func() {

		buf := make([]byte, 1024)
		sTime := time.Now()
		for i := 0; i < 10000; i++ {
			_, err := tap.Write(buf)
			if err != nil {
				t.Fatal(err)
			}
		}
		eTime := time.Now()
		fmt.Println("同步写入10000包耗时 ", eTime.Sub(sTime))
	}()

	func() {
		buf := make([]byte, 1024)
		wg := sync.WaitGroup{}
		wg.Add(10000)
		sTime := time.Now()
		for i := 0; i < 10000; i++ {
			go func() {
				defer wg.Done()
				_, err := tap.Write(buf)
				if err != nil {
					t.Fatal(err)
				}
			}()
		}
		wg.Wait()
		eTime := time.Now()
		fmt.Println("协程写入10000包耗时 ", eTime.Sub(sTime))
	}()
	fmt.Println("end")
}

func TestRead(t *testing.T) {
	fmt.Println("testRead...")
	f := NewTunTapF()
	lAddr := net.IPv4(192, 101, 150, 11)
	network := net.IPv4(190, 101, 150, 0)
	mask := net.IPv4(255, 255, 255, 0)

	tap, err := f.NewTun(lAddr, network, mask)
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	//批量写入内容
	go func() {
		defer wg.Done()
		c, err := net.Dial("udp", "192.101.150.1:456")
		buf := make([]byte, 1000)
		if err != nil {
			t.Fatal(err)
		}

		sTime := time.Now()
		for i := 0; i < 10000; i++ {
			c.Write(buf)
		}
		eTime := time.Now()
		fmt.Println("发送10000 udp 包耗时：", eTime.Sub(sTime))

		winTap := tap.(*tun)

		writeRawPack(winTap.netCfgInstanceId)
	}()

	go func() {
		buf := make([]byte, 2000)
		for i := 0; i < 1000; i++ {
			sTime := time.Now()
			if n, err := tap.Read(buf); err != nil {
				fmt.Println(err)
			} else {
				mem.Put(buf[:n])
			}
			eTime := time.Now()
			fmt.Printf("读取 udp 包耗时：%v\r\n", eTime.Sub(sTime))
			fmt.Println(i)
		}
	}()

	fmt.Println("wait...")
	wg.Wait()
	tap.Close()
	fmt.Println("end")
	time.Sleep(1 * time.Second)
}

func writeRawPack(ifName string) {
	handle, err := pcap.OpenLive(`\Device\NPF_`+ifName, 65536, false, pcap.BlockForever)
	if err != nil {
		log.Fatal("pcap.OpenLive", err)
	}
	buf := make([]byte, 1000)
	for i := 0; i < 10; i++ {
		handle.WritePacketData(buf)
	}
}
