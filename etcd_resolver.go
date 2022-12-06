package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const (
	// 同步实例列表的周期
	syncNSInterval = 5 * time.Second
)

type etcdResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	etcdClient *clientv3.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

func (r *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {
	log.Println("etcd resolver resolve now")
}

func (r *etcdResolver) Close() {
	log.Println("etcd resolver close")
	r.cancel()
}

// 轮询并更新服务的实例
func (r *etcdResolver) watcher() {
	r.updateState()
	ticker := time.NewTicker(syncNSInterval)
	for {
		select {
		case <-ticker.C:
			r.updateState()
		case <-r.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// 更新实例列表
func (r *etcdResolver) updateState() {
	instances := r.getInstances()
	newAddrs := make([]resolver.Address, 0)
	for _, ins := range instances {
		addr := resolver.Address{Addr: fmt.Sprintf("%s:%d", ins.ip, ins.port)}
		// 通过属性存储权重
		addr = SetAddrInfo(addr, AddrInfo{
			Weight: ins.weight,
		})
		newAddrs = append(newAddrs, addr)
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: newAddrs})
}

type Instance struct {
	ip     string
	port   int
	weight int
}

// 获取服务可用的实例
func (r *etcdResolver) getInstances() []*Instance {
	return []*Instance{
		{
			ip:     "127.0.0.1",
			port:   1024,
			weight: 100,
		},
		{
			ip:     "127.0.0.1",
			port:   1025,
			weight: 50,
		},
		{
			ip:     "127.0.0.1",
			port:   1026,
			weight: 25,
		},
	}
}
