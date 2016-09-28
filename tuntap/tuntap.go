package tuntap

import (
	"net"

	"io"
)

// linux 下 tun 类似于文件，默认是程序关闭tun就被删除，每次tun都是新建的。
// windows 下 tun 是虚拟网卡，创建起来需要安装驱动，很慢，而且安装完立刻使用容易出问题。
//         所以一般 windows下网卡是第一次使用时创建，之后就不占删除了。

type TunTapF interface {
	Remove() error
	NewTun(addr net.IP, network net.IP, mask net.IP) (Tun, error)
}

type Tun interface {
	io.Reader
	io.Writer
	io.Closer
	GetName() string
}

type tunTapF struct {
}

func NewTunTapF() TunTapF {
	return &tunTapF{}
}
