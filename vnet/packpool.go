package vnet

const MaxPackSize = 1500

type PackLayer int

const (
	PackLayerEthernet PackLayer = 1
	PackLayerIP       PackLayer = 2
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
