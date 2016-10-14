package ping

import (
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	p, err := Ping("127.0.0.1")
	if err != nil {
		t.Error(err)
	}

	if p > int64((2 * time.Millisecond)) {
		t.Errorf("%v > 2 Millisecond", p)
	}

	p, err = Ping("0.0.0.0")
	if err == nil {
		t.Error("err==nil")
	}
	if p != -1 {
		t.Error("p!=int64Max")
	}

	p, err = Ping("www.baidu.com")
	if err != nil {
		t.Error(err)
	}

	if p > int64(2*time.Second) || p < int64(1*time.Millisecond) {
		t.Error("p>2*time.Second||p<1*time.Millisecond")
	}

}

func TestFindTime(t *testing.T) {
	i, err := findTime(`PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.035 ms

--- 127.0.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.035/0.035/0.035/0.000 ms
`)
	if err != nil {
		t.Error(err)
	}

	if i != int64(0.035*float64(time.Millisecond)) {
		t.Errorf("%v != 0.035", i)
	}

	i, err = findTime(`PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=100 ms

--- 127.0.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.035/0.035/0.035/0.000 ms
`)
	if err != nil {
		t.Error(err)
	}

	if i != int64(100*float64(time.Millisecond)) {
		t.Errorf("%v != 100", i)
	}

	i, err = findTime(`正在 Ping 127.0.0.1 具有 32 字节的数据:
来自 127.0.0.1 的回复: 字节=32 时间<1ms TTL=128

127.0.0.1 的 Ping 统计信息:
    数据包: 已发送 = 1，已接收 = 1，丢失 = 0 (0% 丢失)，
往返行程的估计时间(以毫秒为单位):
    最短 = 0ms，最长 = 0ms，平均 = 0ms`)
	if err != nil {
		t.Error(err)
	}

	if i != int64(1*float64(time.Millisecond)) {
		t.Errorf("%v != 1", i)
	}
}
