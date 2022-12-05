package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const (
	// 同步实例列表的周期
	syncNSInterval = time.Second
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
		newAddrs = append(newAddrs, resolver.Address{Addr: ins})
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: newAddrs})
}

// 获取服务可用的实例
func (r *etcdResolver) getInstances() []string {
	// todo: 从 etcd 拉取服务实例
	return []string{"127.0.0.1:1024", "127.0.0.1:1025"}
}
