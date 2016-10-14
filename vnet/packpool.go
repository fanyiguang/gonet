package vnet

import (
	"encoding/binary"
	"fmt"
	"net"
)

const MaxPackSize = 2000

type PackLayer int

const (
	PackLayerLink        PackLayer = 0x10 //以太网
	PackLayerNetwork     PackLayer = 0x20 // ip
	PackLayerTransport   PackLayer = 0x30 //tcp udp
	PackLayerApplication PackLayer = 0x40 //
)

type Pack interface {

	// 所有数据
	// 为了减少内存拷贝，会预先保留二次封装包包头的空间，使得封包时不用再次拷贝。
	GetAllData() []byte

	// tun/tap 包数据
	GetPrefixData() []byte
	GetPackData() []byte
	SetPackLength(int)

	// 释放本包占用的内存
	Free()

	Reset()

	GetPackLayer() PackLayer

	// 不安全，非拷贝副本
	GetIP() (net.IP, net.IP, error)
	GetPort() (int, int, error)
}

type Freer interface {
	Free(Pack)
}

type PackPool interface {
	Malloc() Pack
}

type pack struct {
	data         []byte
	prefixSize   int
	packDataSize int
	free         Freer
	packLayer    PackLayer
}

type packPool struct {
	// 之前测试过，多线程时还是 chan 实现池效果最好
	pool chan Pack
	new  func() Pack
}

func (p *pack) GetPrefixData() []byte {
	return p.data[:p.prefixSize]
}
func (p *pack) GetAllData() []byte {
	return p.data[:p.prefixSize+p.packDataSize]
}

func (p *pack) GetPackData() []byte {
	return p.data[p.prefixSize : p.prefixSize+p.packDataSize]
}

func (p *pack) SetPackLength(size int) {
	p.packDataSize = size
}

func (p *pack) Free() {
	if p.free != nil {
		p.free.Free(p)
	}
}

func (p *pack) GetPackLayer() PackLayer {
	return p.packLayer
}

func (p *pack) Reset() {
	p.packDataSize = len(p.data) - p.prefixSize
	pd := p.GetPrefixData()
	for i, _ := range pd {
		pd[i] = 0
	}
}

func (p *pack) GetIP() (sip, oip net.IP, rerr error) {
	data, err := p.GetNetworkLayerData()
	if err != nil {
		return nil, nil, err
	}

	switch data[0] & 0xf0 {
	case 0x40:
		// ipv4
		if len(data) < 20 {
			return nil, nil, fmt.Errorf("长度不正确")
		}
		return net.IP(data[12:16]), net.IP(data[16:20]), nil
	case 0x60:
		// ipv6 忽略
		return nil, nil, fmt.Errorf("ipv6 未实现")
	default:
		return nil, nil, fmt.Errorf("未知协议 0x%X", data[0])
	}
}

func (p *pack) GetNetworkLayerData() ([]byte, error) {
	switch p.GetPackLayer() {
	case PackLayerNetwork:
		return p.GetPackData(), nil
	default:
		return nil, fmt.Errorf("不支持")
	}
}

func (p *pack) GetTransportLayerData() (int, []byte, error) {
	data := p.GetPackData()
	packLayer := p.GetPackLayer()
	if packLayer == PackLayerTransport {
		return 0, data, nil
	}

	if packLayer < PackLayerNetwork {
		ldata, err := p.GetNetworkLayerData()
		if err != nil {
			return 0, nil, err
		}
		data = ldata
	} else if p.GetPackLayer() >= PackLayerApplication {
		return 0, nil, fmt.Errorf("不支持")
	}

	switch data[0] & 0xf0 {
	case 0x40:
		// ipv4
		if len(data) < 20 {
			return 0, nil, fmt.Errorf("长度太短。")
		}
		headLen := int(data[0]&0x0f) * 4
		allLen := binary.BigEndian.Uint16(data[2:4])
		ptype := data[9]

		if len(data) < int(allLen) {
			return 0, nil, fmt.Errorf("长度太短")
		}
		return int(ptype), data[headLen:allLen], nil
	case 0x60:
		// ipv6 忽略
		return 0, nil, fmt.Errorf("ipv6 未实现")
	default:
		return 0, nil, fmt.Errorf("未知协议 0x%X", data[0])
	}
}

func (p *pack) GetPort() (sport, oport int, rerr error) {
	ptype, data, err := p.GetTransportLayerData()
	if err != nil {
		return 0, 0, err
	}
	fmt.Println(data)

	switch ptype {
	case 17:
		//udp
		if len(data) < 4 {
			return 0, 0, fmt.Errorf("长度不正确")
		}
		return int(binary.BigEndian.Uint16(data[0:2])), int(binary.BigEndian.Uint16(data[2:4])), nil

	case 6:
		//tcp
		if len(data) < 4 {
			return 0, 0, fmt.Errorf("长度不正确")
		}
		return int(binary.BigEndian.Uint16(data[0:2])), int(binary.BigEndian.Uint16(data[2:4])), nil
	default:
		return 0, 0, fmt.Errorf("未知 Transport 协议 %v", ptype)
	}
}

func newPack(data []byte, prefixSize int, packLayer PackLayer, free Freer) Pack {
	p := &pack{
		data:         data,
		prefixSize:   prefixSize,
		packDataSize: len(data) - prefixSize,
		packLayer:    packLayer,
		free:         free,
	}
	return p
}

func NewPool(prefixSize, MaxPackDataSize int, packLayer PackLayer) PackPool {
	pool := &packPool{
		pool: make(chan Pack, 100),
	}

	pool.new = func() Pack {
		return newPack(make([]byte, prefixSize+MaxPackDataSize), prefixSize, packLayer, pool)
	}

	return pool
}

func (pool *packPool) Malloc() Pack {
	var p Pack
	select {
	case p = <-pool.pool:
	default:
		p = pool.new()
	}
	p.Reset()
	return p
}

func (pool *packPool) Free(p Pack) {
	select {
	case pool.pool <- p:
	default:
		return
	}
}
