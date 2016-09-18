package vnet

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestPack(t *testing.T) {
	pool := NewPool(10, MaxPackSize, PackLayerIP)
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
	pool := NewPool(10, MaxPackSize, PackLayerIP)

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
