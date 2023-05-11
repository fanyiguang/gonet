//go:build windows
// +build windows

package route

import (
	"bytes"
	"net"
	"testing"
)

func TestRouteTableRoutes(t *testing.T) {
	rt := NewRouteTable()
	rT := rt.(*routeTable)
	rs, err := rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}

	var defaultRoute *RouteRow
	for _, r := range rs {
		if r.ForwardType == 4 && r.ForwardProto == 3 &&
			bytes.Equal([]byte{0, 0, 0, 0}, []byte(r.GetForwardDest())) &&
			bytes.Equal([]byte{0, 0, 0, 0}, []byte(r.GetForwardMask())) {
			defaultRoute = &r
		}
	}

	if defaultRoute == nil {
		t.Fatal("defaultRoute == nil")
	}

	if bytes.Equal([]byte(defaultRoute.GetForwardNextHop()), []byte{0, 0, 0, 0}) {
		t.Fatal("defaultRoute.ForwardNextHop==[]byte{0,0,0,0}")
	}
}

func TestAddRouteAndDelRoute(t *testing.T) {
	rt := NewRouteTable()
	rT := rt.(*routeTable)

	if err := rt.ResetRoute(); err != nil {
		t.Errorf("复位路由表失败，%v", err)
	}

	rs, err := rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}

	var defaultRoute *RouteRow
	for _, r := range rs {
		if r.ForwardType == 4 && r.ForwardProto == 3 &&
			bytes.Equal([]byte{0, 0, 0, 0}, r.ForwardDest[:]) &&
			bytes.Equal([]byte{0, 0, 0, 0}, r.ForwardMask[:]) {
			t := r
			defaultRoute = &t
		}
	}

	if defaultRoute == nil {
		t.Fatal("defaultRoute == nil")
	}

	r := *defaultRoute

	// 测试增加

	// 准备
	if n := copy(r.ForwardDest[:], []byte{222, 111, 123, 0}); n != 4 {
		t.Fatal(n, "!=4")
	}
	if n := copy(r.ForwardMask[:], []byte{255, 255, 255, 0}); n != 4 {
		t.Fatal(n, "!=4")
	}
	r.ForwardMetric1 = routeMetric
	r.ForwardNextHop[3] += 1

	// 增加
	if err := rT.addRoute(&r); err != nil {
		t.Fatal(err)
	}

	// 检查是否者增加成功
	rs, err = rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, rr := range rs {
		if bytes.Equal(rr.ForwardDest[:], r.ForwardDest[:]) && bytes.Equal(rr.ForwardMask[:], r.ForwardMask[:]) &&
			bytes.Equal(rr.ForwardNextHop[:], r.ForwardNextHop[:]) {
			count++
		}
	}
	if count != 1 {
		t.Errorf("%v != 1", count)
	}

	// 测试删除
	if err := rT.delRoute(&r); err != nil {
		t.Error(err)
	}

	//测试是否已经不存在
	rs, err = rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}
	count = 0
	for _, rr := range rs {
		if bytes.Equal(rr.ForwardDest[:], r.ForwardDest[:]) && bytes.Equal(rr.ForwardMask[:], r.ForwardMask[:]) &&
			bytes.Equal(rr.ForwardNextHop[:], r.ForwardNextHop[:]) && rr.ForwardMetric2 == routeMetric {
			count++
		}
	}
	if count != 0 {
		t.Errorf("%v != 0", count)
	}

}

func TestAddNetRoutesAndResetRoutes(t *testing.T) {
	rt := NewRouteTable()
	rT := rt.(*routeTable)

	if err := rt.ResetRoute(); err != nil {
		t.Errorf("复位路由表失败，%v", err)
	}

	rs, err := rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}

	var defaultRoute *RouteRow
	for _, r := range rs {
		if r.ForwardType == 4 && r.ForwardProto == 3 &&
			bytes.Equal([]byte{0, 0, 0, 0}, r.ForwardDest[:]) &&
			bytes.Equal([]byte{0, 0, 0, 0}, r.ForwardMask[:]) {
			t := r
			defaultRoute = &t
		}
	}
	if defaultRoute == nil {
		t.Fatal("无默认网关，无法测试 AddNetRoutes")
	}

	routes := []Route{
		{
			Network: net.IPv4(142, 233, 111, 0),
			Mask:    24,
		},
		{
			Network: net.IPv4(142, 233, 222, 0),
			Mask:    24,
		},
	}

	if err := rt.AddNetRoutes(routes); err != nil {
		t.Fatal(err)
	}

	rs, err = rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for _, r := range rs {
		if r.GetForwardDest().Equal(net.IPv4(142, 233, 111, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			if r.GetForwardNextHop().Equal(defaultRoute.GetForwardNextHop()) == false {
				t.Fatal(defaultRoute.GetForwardNextHop(), "!=", r.GetForwardNextHop())
			}
			count++
		}
		if r.GetForwardDest().Equal(net.IPv4(142, 233, 222, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			if r.GetForwardNextHop().Equal(defaultRoute.GetForwardNextHop()) == false {
				t.Fatal(defaultRoute.GetForwardNextHop(), "!=", r.GetForwardNextHop())
			}
			count++
		}
	}

	if count != 2 {
		t.Fatal(count, "!=2")
	}

	if err := rt.ResetRoute(); err != nil {
		t.Fatal(err)
	}

	rs, err = rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	count = 0
	for _, r := range rs {
		if r.GetForwardDest().Equal(net.IPv4(142, 233, 111, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			count++
		}
		if r.GetForwardDest().Equal(net.IPv4(142, 233, 222, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			count++
		}
	}
	if count != 0 {
		t.Fatal(count, "!=0")
	}
}

func TestAddVpnRoutesAndResetRoutes(t *testing.T) {
	rt := NewRouteTable()
	rT := rt.(*routeTable)

	if err := rt.ResetRoute(); err != nil {
		t.Errorf("复位路由表失败，%v", err)
	}

	rs, err := rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}

	routes := []Route{
		{
			Network: net.IPv4(183, 233, 111, 0),
			Mask:    24,
		},
		{
			Network: net.IPv4(183, 233, 222, 0),
			Mask:    24,
		},
		{
			Network: net.IPv4(0, 0, 0, 0),
			Mask:    0,
		},
	}

	if err := rt.AddVpnRoutes(routes, net.IPv4(192, 168, 16, 0), net.IPv4(255, 255, 255, 0), net.IPv4(192, 168, 16, 11)); err != nil {
		t.Fatal(err)
	}

	rs, err = rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for _, r := range rs {
		if r.GetForwardDest().Equal(net.IPv4(183, 233, 111, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			if r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) == false {
				t.Error(net.IPv4(192, 168, 16, 11), "!=", r.GetForwardNextHop())
			}
			count++
		}
		if r.GetForwardDest().Equal(net.IPv4(183, 233, 222, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			if r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) == false {
				t.Error(net.IPv4(192, 168, 16, 11), "!=", r.GetForwardNextHop())
			}
			count++
		}
		if r.GetForwardDest().Equal(net.IPv4(0, 0, 0, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(128, 0, 0, 0)) {
			if r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) == false {
				t.Error(net.IPv4(192, 168, 16, 11), "!=", r.GetForwardNextHop())
			}
			count++
		}
		if r.GetForwardDest().Equal(net.IPv4(128, 0, 0, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(128, 0, 0, 0)) {
			if r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) == false {
				t.Error(net.IPv4(192, 168, 16, 11), "!=", r.GetForwardNextHop())
			}
			count++
		}
		if r.GetForwardDest().Equal(net.IPv4(0, 0, 0, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(0, 0, 0, 0)) && r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) {
			t.Errorf("默认路由没有特殊处理！ %#v", r)
		}
	}

	if count != 4 {
		t.Error(count, "!=4")
	}

	if err := rt.ResetRoute(); err != nil {
		t.Error(err)
	}

	rs, err = rT.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range rs {
		if r.GetForwardDest().Equal(net.IPv4(183, 233, 111, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			t.Errorf("%v 路由没有删除。%#v", net.IPv4(123, 233, 111, 0), r)
		}
		if r.GetForwardDest().Equal(net.IPv4(183, 233, 222, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(255, 255, 255, 0)) {
			t.Errorf("%v 路由没有删除。%#v", net.IPv4(123, 233, 222, 0), r)
			count++
		}
		if r.GetForwardDest().Equal(net.IPv4(0, 0, 0, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(128, 0, 0, 0)) {
			if r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) {
				t.Errorf("%v 路由没有删除。%#v", net.IPv4(0, 0, 0, 0), r)
				count++
			}
		}
		if r.GetForwardDest().Equal(net.IPv4(128, 0, 0, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(128, 0, 0, 0)) {
			if r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) {
				t.Errorf("%v 路由没有删除。%#v", net.IPv4(128, 0, 0, 0), r)
				count++
			}
		}
		if r.GetForwardDest().Equal(net.IPv4(0, 0, 0, 0)) &&
			r.GetForwardMask().Equal(net.IPv4(0, 0, 0, 0)) && r.GetForwardNextHop().Equal(net.IPv4(192, 168, 16, 11)) {
			t.Errorf("%v 路由没有删除。%#v", net.IPv4(0, 0, 0, 0), r)
			count++
		}
	}

}
