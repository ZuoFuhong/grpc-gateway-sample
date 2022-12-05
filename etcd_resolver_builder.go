package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

type EtcdResolverBuilder struct {
	etcdClient *clientv3.Client
}

func NewEtcdResolverBuilder() *EtcdResolverBuilder {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println("client get etcd failed,error", err)
		panic(err)
	}
	return &EtcdResolverBuilder{
		etcdClient: etcdClient,
	}
}

func (erb *EtcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &etcdResolver{
		target:     target,
		cc:         cc,
		etcdClient: erb.etcdClient,
		ctx:        ctx,
		cancel:     cancel,
	}
	// 启动协程
	go r.watcher()
	return r, nil
}

func (erb *EtcdResolverBuilder) Scheme() string {
	return "etcd"
}
