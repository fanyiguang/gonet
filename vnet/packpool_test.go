package vnet

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestPack(t *testing.T) {
	pool := NewPool(10, MaxPackSize, PackLayerNetwork)
	p := pool.Malloc()
	if len(p.GetPrefixData()) != 10 {
		t.Error("len(p.GetPrefixData())=", len(p.GetPrefixData()))
	}
	if len(p.GetPackData()) != MaxPackSize {
		t.Error("len(p.GetPackData())=", len(p.GetPackData()))
	}
	if len(p.GetAllData()) != 10+MaxPackSize {
		t.Error("len(p.GetAllData())=", len(p.GetAllData()))
	}

	pd := p.GetPackData()
	for i, _ := range pd {
		pd[i] = 123
	}

	for _, v := range p.GetPrefixData() {
		if v == 123 {
			t.Error(v)
		}
	}

	pd = p.GetPrefixData()
	for i, _ := range pd {
		pd[i] = 111
	}

	p.SetPackLength(110)

	if len(p.GetPrefixData()) != 10 {
		t.Error("len(p.GetPrefixData())=", len(p.GetPrefixData()))
	}
	if len(p.GetPackData()) != 110 {
		t.Error("len(p.GetPackData())=", len(p.GetPackData()))
	}
	if len(p.GetAllData()) != 10+110 {
		t.Error("len(p.GetAllData())=", len(p.GetAllData()))
	}

	p.Free()
	p2 := pool.Malloc()
	for _, v := range p2.GetPrefixData() {
		if v != 0 {
			t.Error(v)
		}
	}

	if !bytes.Equal(p2.GetPackData(), p.GetPackData()) {
		t.Error("p!=p2")
	}
	if len(p.GetPrefixData()) != 10 {
		t.Error("len(p.GetPrefixData())=", len(p.GetPrefixData()))
	}
	if len(p.GetPackData()) != MaxPackSize {
		t.Error("len(p.GetPackData())=", len(p.GetPackData()))
	}
	if len(p.GetAllData()) != 10+MaxPackSize {
		t.Error("len(p.GetAllData())=", len(p.GetAllData()))
	}
}

func TestPool(t *testing.T) {
	pool := NewPool(10, MaxPackSize, PackLayerNetwork)

	func() {
		sTime := time.Now()
		for i := 0; i < 1000000; i++ {
			_ = pool.Malloc()
		}
		fmt.Println("单独申请执行耗时：", time.Now().Sub(sTime))
	}()

	func() {
		sTime := time.Now()
		for i := 0; i < 1000000; i++ {
			p := pool.Malloc()
			p.Free()
		}
		fmt.Println("申请+释放执行耗时：", time.Now().Sub(sTime))
	}()
}

/*

单独申请执行耗时： 713.287ms
申请+释放执行耗时： 121.6307ms

*/

func ToBytes(d string) ([]byte, error) {
	ss := strings.Split(d, " ")
	r := make([]byte, 0, len(d)/3)
	for _, s := range ss {
		i, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析错误，%v", err)
		}

		r = append(r, byte(i))
	}
	return r, nil
}

func ToPack(s string, t PackLayer) (Pack, error) {

	d1, err := ToBytes(s)
	if err != nil {
		return nil, err
	}

	p := newPack(make([]byte, 2048), 3, t, nil)
	copy(p.GetPackData(), d1)
	p.SetPackLength(len(d1))

	return p, nil
}

func TestPack_TCP(t *testing.T) {
	s1 := `45 00 00 28 b8 b5 00 00 2d 06 33 23 dc ff 02 99 c0 a8 01 b7 01 bb cb c5 74 f7 b1 66 60 65 91 05 50 11 00 e5 27 ad 00 00 00 00 00 00 00 00`

	p, err := ToPack(s1, PackLayerNetwork)
	if err != nil {
		t.Fatal(err)
	}

	if p.GetPackLayer() != PackLayerNetwork {
		t.Fatal("!=PackLayerNetwork")
	}

	sip, oip, err := p.GetIP()
	if !sip.Equal(net.IPv4(220, 255, 2, 153)) || !oip.Equal(net.IPv4(192, 168, 1, 183)) {
		t.Error(sip, oip, err)
	}

	sport, oport, err := p.GetPort()
	if sport != 443 || oport != 52165 {
		t.Error(sport, oport, err)
	}

}

func TestPack_UDP(t *testing.T) {
	s1 := `45 00 00 20 67 a5 00 00 80 11 eb 50 c0 a8 01 b7 2d 20 f8 57 27 09 27 09 00 0c 30 55 4b cc 4d cb`

	p, err := ToPack(s1, PackLayerNetwork)
	if err != nil {
		t.Fatal(err)
	}

	if p.GetPackLayer() != PackLayerNetwork {
		t.Fatal("!=PackLayerNetwork")
	}

	sip, oip, err := p.GetIP()
	if !sip.Equal(net.IPv4(192, 168, 1, 183)) || !oip.Equal(net.IPv4(45, 32, 248, 87)) {
		t.Error(sip, oip, err)
	}

	sport, oport, err := p.GetPort()
	if sport != 9993 || oport != 9993 {
		t.Error(sport, oport, err)
	}

}

func TestPack_DNS(t *testing.T) {
	s1 := `45 00 00 3e 42 cc 00 00 80 11 73 da c0 a8 01 b7 c0 a8 01 01 e8 94 00 35 00 2a 84 8e 17 3c 01 00 00 01 00 00 00 00 00 00 05 69 6e 62 6f 78 06 67 6f 6f 67 6c 65 03 63 6f 6d 00 00 01 00 01`

	p, err := ToPack(s1, PackLayerNetwork)
	if err != nil {
		t.Fatal(err)
	}

	if p.GetPackLayer() != PackLayerNetwork {
		t.Fatal("!=PackLayerNetwork")
	}

	sip, oip, err := p.GetIP()
	if !sip.Equal(net.IPv4(192, 168, 1, 183)) || !oip.Equal(net.IPv4(192, 168, 1, 1)) {
		t.Error(sip, oip, err)
	}

	sport, oport, err := p.GetPort()
	if sport != 59540 || oport != 53 {
		t.Error(sport, oport, err)
	}

}
