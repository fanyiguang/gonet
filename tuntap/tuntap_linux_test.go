package tuntap

import (
	"net"
	"testing"
)

func TestTunAll(t *testing.T) {

	tf := NewTunTapF()

	ip := net.IPv4(192, 168, 194, 31)

	_, err := tf.NewTun(ip, net.IPv4(192, 168, 194, 0), net.IPv4(255, 255, 255, 0))
	if err != nil {
		t.Fatal(err)
	}

	ifs, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range ifs {
		addrs, err := i.Addrs()
		if err != nil {
			t.Errorf(err)
		}

		for _, addr := range addrs {
			if ip.String() == addr.String() {
				return
			}
		}
	}
	t.Fatal("未找到创建的网卡。")

	//fmt.Println("=======================")
	//time.Sleep(100 * time.Second)
	/*
		c := make(chan error)

		go func() {
			for {
				data, err := tun.ReadPack()
				if err != nil {
					c <- err
					return
				}

				defer mem.Put(data)

				switch data[0] & 0xf0 {
				case 0x40:
					packet := gopacket.NewPacket(data, layers.LayerTypeIPv4, gopacket.Default)
					fmt.Println(packet.String())
					c <- nil
				case 0x60:
					packet := gopacket.NewPacket(data, layers.LayerTypeIPv6, gopacket.Default)
					fmt.Println(packet.String())
					c <- nil
				}
			}
		}()

		func() {
				buf := make([]byte, 1024)
				err := tun.WritePack(buf)
				if err != nil {
					t.Fatal(err)
				}

		}()

		cmd := exec.Command("ping", "-c", "3", "192.168.184.100")
		out, err := cmd.Output()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(out))

		timeoutchan := time.After(5 * time.Second)

		select {

		case err := <-c:
			if err != nil {
				t.Fatal(err)
			}
			return
		case <-timeoutchan:
			t.Fatal("timeout")
		}*/

}
