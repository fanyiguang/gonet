package gprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// 注意，长度不足时会直接 panic: runtime error: index out of range
func ReadUint8(b *[]byte) uint8 {
	i := (*b)[0]
	*b = (*b)[1:]
	return i
}
func ReadUint16(b *[]byte) uint16 {
	i := binary.BigEndian.Uint16(*b)
	*b = (*b)[2:]
	return i
}
func ReadUint32(b *[]byte) uint32 {
	i := binary.BigEndian.Uint32(*b)
	*b = (*b)[4:]
	return i
}
func ReadUint64(b *[]byte) uint64 {
	i := binary.BigEndian.Uint64(*b)
	*b = (*b)[8:]
	return i
}

func WriteUint8(b *[]byte, v uint8) {
	*b = append(*b, v)
}
func WriteUint16(b *[]byte, v uint16) {
	var bbuf [2]byte
	buf := bbuf[:]
	binary.BigEndian.PutUint16(buf, v)
	*b = append(*b, buf...)
}
func WriteUint32(b *[]byte, v uint32) {
	var bbuf [4]byte
	buf := bbuf[:]
	binary.BigEndian.PutUint32(buf, v)
	*b = append(*b, buf...)
}
func WriteUint64(b *[]byte, v uint64) {
	var bbuf [8]byte
	buf := bbuf[:]
	binary.BigEndian.PutUint64(buf, v)
	*b = append(*b, buf...)
}
func WriteInt32(b *[]byte, v int32) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, v); err != nil {
		panic(err)
	}
	*b = append(*b, buf.Bytes()...)
}
func WriteInt64(b *[]byte, v int64) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, v); err != nil {
		panic(err)
	}
	*b = append(*b, buf.Bytes()...)
}
func ReadInt32(b *[]byte) int32 {
	buf := bytes.NewReader((*b)[:4])
	i := int32(0)
	if err := binary.Read(buf, binary.BigEndian, &i); err != nil {
		panic(err)
	}
	*b = (*b)[4:]
	return i
}
func ReadInt64(b *[]byte) int64 {
	buf := bytes.NewReader((*b)[:8])
	i := int64(0)
	if err := binary.Read(buf, binary.BigEndian, &i); err != nil {
		panic(err)
	}
	*b = (*b)[8:]
	return i
}

func Read2String(b *[]byte) (s string) {

	l := ReadUint16(b)
	if int(l) > len(*b) {
		panic(fmt.Errorf("长度不足。"))
	}
	s = string((*b)[:l])
	*b = (*b)[l:]
	return s
}

func Write2String(b *[]byte, v string) {
	if len(v) > 0xFFFF {
		panic("字符串超出允许长度。")
	}
	WriteUint16(b, uint16(len(v)))
	*b = append(*b, []byte(v)...)
}

//警告，直接返回的 b 的切片，如果b会被修改请自己拷贝一份副本。
func Read2BytesUsafe(b *[]byte) (s []byte) {

	l := ReadUint16(b)
	if int(l) > len(*b) {
		panic(fmt.Errorf("长度不足。"))
	}
	s = (*b)[:l]
	*b = (*b)[l:]
	return s
}
func Read2Bytes(b *[]byte) (s []byte) {
	ob := Read2BytesUsafe(b)
	nb := make([]byte, len(ob))
	copy(nb, ob)
	return nb
}

func Write2Bytes(b *[]byte, v []byte) {
	if len(v) > 0xFFFF {
		panic("字符串超出允许长度。")
	}
	WriteUint16(b, uint16(len(v)))
	*b = append(*b, v...)
}
