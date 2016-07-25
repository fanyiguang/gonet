// +build !windows

package route

import (
	"fmt"
	"net"

)

type routeTable struct {
}

func NewRouteTable() RouteTable {
	return &routeTable{}
}

func (rt *routeTable) AddNetRoutes(routes []Route) error {
	return fmt.Errorf("当前系统不支持此功能。")
}
func (rt *routeTable) AddVpnRoutes(routes []Route, network, mask, gIp net.IP) error {
	return fmt.Errorf("当前系统不支持此功能。")
}

func (rt *routeTable) ResetRoute() error {
	return fmt.Errorf("当前系统不支持此功能。")
}
