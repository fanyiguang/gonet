package nettool

import (
	"fmt"
	"time"
)

type SetLingerer interface {
	SetLinger(sec int) error
}
type SetNoDelayer interface {
	SetNoDelay(noDelay bool) error
}

type SetReadBufferer interface {
	SetReadBuffer(bytes int) error
}
type SetWriteBufferer interface {
	SetWriteBuffer(bytes int) error
}
type SetReadBuffer interface {
	SetReadBuffer(bytes int) error
}

type SetReadDeadlineer interface {
	SetReadDeadline(t time.Time) error
}

type SetWriteDeadlineer interface {
	SetWriteDeadline(t time.Time) error
}

//SetDeadline 包含 net.Conn

// 开关 Delay 算法
// noDelay = true 关闭
// false 性能更好，但是延迟高
func SetNoDelay(conn interface{}, noDelay bool) error {
	ccd, _ := conn.(SetNoDelayer)
	if ccd == nil {
		return fmt.Errorf("conn 未提供 SetNoDelay 方法。")
	}
	ccd.SetNoDelay(noDelay)
	return nil
}

func SetLinger(c interface{}, sec int) error {
	ccd, _ := c.(SetLingerer)
	if ccd == nil {
		return fmt.Errorf("conn 未提供 SetNoDelay 方法。")
	}
	ccd.SetLinger(sec)
	return nil
}

func SetSetReadDeadline(c interface{}, t time.Time) error {
	ccd, _ := c.(SetReadDeadlineer)
	if ccd == nil {
		return fmt.Errorf("conn 未提供 SetReadDeadline 方法。")
	}
	ccd.SetReadDeadline(t)
	return nil
}

func SetWriteDeadline(c interface{}, t time.Time) error {
	ccd, _ := c.(SetWriteDeadlineer)
	if ccd == nil {
		return fmt.Errorf("conn 未提供 SetWriteDeadline 方法。")
	}
	ccd.SetWriteDeadline(t)
	return nil
}
