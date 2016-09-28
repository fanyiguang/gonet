package vnet

import "sync"

// 虚拟网络
// 虚拟网络允许接入多个来源，负责将对应的请求转发到对应的ip协议栈
// 这里相当于 ip 路由器
// 目前只允许单个ip接入，后期计划支持子网加入。
//
// 尽可能实现内存 0 拷贝，
//
type vnet struct {
	rwm sync.RWMutex
}
