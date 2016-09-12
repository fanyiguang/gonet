package tuntap

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

const (
	IFF_NO_PI = 0x10
	IFF_TUN   = 0x01
	IFF_TAP   = 0x02
	TUNSETIFF = 0x400454CA
)

type tun struct {
	mtu        int
	actualName string
	fd         *os.File
}

const (
	flagTruncated = 0x1

	iffTun      = 0x1
	iffTap      = 0x2
	iffOneQueue = 0x2000
	iffnopi     = 0x1000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func (tf *tunTapF) NewTun(addr net.IP, network net.IP, mask net.IP) (Tun, error) {
	fpath := "/dev/net/tun"
	ifPattern := ""

	file, err := os.OpenFile(fpath, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("打开 tun 设备 %v 失败,%v", fpath, err)
	}

	var req ifReq
	req.Flags = 0
	copy(req.Name[:15], ifPattern)

	req.Flags |= iffTun
	req.Flags |= iffnopi

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if r1 != 0 {
		return nil, err
	}

	t := tun{
		actualName: strings.Trim(string(req.Name[:]), "\000"),
		fd:         file,
	}

	err = t.setupAddress(addr.String(), mask.String())
	if err != nil {
		return nil, fmt.Errorf("设置IP错误,%v", err)
	}

	return &t, nil
}

func (t *tunTapF) Remove() error {
	return nil
}

func (t *tun) setupAddress(addr, mask string) error {
	cmd := exec.Command("ifconfig", t.actualName, addr,
		"netmask", mask, "mtu", "1500")
	log.Printf("[DEBG] ifconfig command: %v", strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		log.Printf("[EROR] Linux ifconfig failed: %v.", err)
		return err
	}
	return nil
}

func (t *tun) Close() error {
	return t.Close()
}

func (t *tun) Write(buf []byte) (int, error) {
	return t.fd.Write(buf)
}

func (t *tun) Read(buf []byte) (int, error) {
	return t.fd.Read(buf)
}

func (t *tun) SetDns(ip []net.IP) error {
	return fmt.Errorf("linux 系统不支持。")
}

func (t *tun) GetName() string {
	return t.actualName
}
