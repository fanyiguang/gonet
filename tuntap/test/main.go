package main

import (
	"net"
	"time"

	"fmt"

	"github.com/gamexg/gonet/tuntap"
)

// 尝试手工实现 tcp 协议，这样测试一下性能。

func main() {
	f := tuntap.NewTunTapF()
	lAddr := net.IPv4(192, 168, 111, 1)
	network := net.IPv4(192, 168, 111, 0)
	mask := net.IPv4(255, 255, 255, 0)

	tap, err := f.NewTun(lAddr, network, mask)
	if err != nil {
		panic(err)
	}
	defer tap.Close()

	go func() {
		b := make([]byte, 2048)
		for {
			_, err := tap.Read(b)
			if err != nil {
				panic(err)
			}
		}
	}()

	func() {
		buf := make([]byte, 1024)
		sTime := time.Now()
		for i := 0; i < 1000000; i++ {
			_, err := tap.Write(buf)
			if err != nil {
				panic(err)
			}
		}
		eTime := time.Now()
		fmt.Println("同步写入1000000包耗时：", eTime.Sub(sTime))
		fmt.Println("平均耗时： ", eTime.Sub(sTime)/1000000)
	}()

	/*
		这个是标准异步 WriteFile 速度：
			同步写入1000000包耗时： 8.9125098s
			平均耗时：  8.912µs

		另外本来应该是高速的完成端口模型，速度可怜：
			同步写入10000包耗时： 13.7945879s
			平均耗时：  1.379458ms

		看起来可能驱动实现有问题造成完成端口基本没法工作。
		悲剧的我还修复了微软完成端口库无法在 windows xp 下工作的问题。

		另外还有一个选择，即完全阻塞 read 、write 模式。
		但是 windows 下同步读写是串行模式，即同一时间只允许一个操作，读会阻塞写。
		参考：https://msdn.microsoft.com/en-us/library/windows/desktop/aa363858(v=vs.85).aspx
		If this flag is specified, the file can be used for simultaneous read and write operations.
		If this flag is not specified, then I/O operations are serialized, even if the calls to the read and write functions specify an OVERLAPPED structure.


	*/

	time.Sleep(1 * time.Hour)
}

/*

tcp 协议要求使用4个元素确定一个连接，虽然一些系统一个端口只能建立一个对外连接连接，但是这里还是按照 tcp 标准来实现吧。

计划只模拟远端设备，并不管本地的ip等。

计划结构是这样的：

远端池->远端ip设备->远端端口(连接池)->单条连接

这里狗屎依靠tcp协议的4个元素进行的分路。

最后一步需要注意下，计划使用窗口大小而不是丢包来限制客户端的发包速度。




*/
