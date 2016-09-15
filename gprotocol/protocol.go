package gprotocol

import (
	"encoding/binary"
	"io"

	"bitbucket.org/jack/jackproxy/common/mem"
)

/*
自身的一些协议

自身的协议可以考虑的比较周全些

现在有2个问题，
接收端可能不清楚当前是什么命令

*/

func WritePackLength(pack []byte) {
	l := len(pack) - 3
	if l < 0 {
		panic("包长度小于3，包头尺寸不正确。")
	}

	if l > 0xFFFF {
		panic("包长度太长。")
	}

	binary.BigEndian.PutUint16(pack[1:3], uint16(l))
}

func ReadPackLength(r io.Reader) (uint16, error) {
	var lb [2]byte
	if _, err := io.ReadFull(r, lb[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(lb[:]), nil
}

func ReadPack(r io.Reader) ([]byte, error) {
	l, err := ReadPackLength(r)
	if err != nil {
		return nil, err
	}
	buf := mem.Get(int(l))

	if n, err := io.ReadFull(r, buf); err != nil {
		return buf[:n], err
	}

	return buf, nil
}

func FreePack(buf []byte) {
	mem.Put(buf)
}

func ReadPType(r io.Reader) (byte, error) {
	var ptype byte
	if err := binary.Read(r, binary.BigEndian, &ptype); err != nil {
		return 0, err
	}
	return ptype, nil
}
