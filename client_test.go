package main

import (
	"context"
	"fmt"
	pb "github.com/ZuoFuhong/grpc_gateway_best_practices/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"net/url"
	"testing"
	"time"
)

func Test_Client(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:1024", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	stub := pb.NewGoWalletManageSvrClient(conn)
	rpcRsp, err := stub.CreateWallet(context.Background(), &pb.CreateWalletReq{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rpcRsp.GetAddress())
}

func Test_ParseTarget(t *testing.T) {
	u, err := url.Parse("etcd://Test/go_wallet_manage_svr")
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Scheme, u.Host, u.Path, "opaque:", u.Opaque)
}

func Test_NameResolver(t *testing.T) {
	// 注册服务发现
	b := NewEtcdResolverBuilder()
	resolver.Register(b)
	// gRPC 提供两种负载均衡策略 pick_first、round_robin, 默认的策略 pick_first
	// 自定义实现 "加权轮询" 负载策略：weighted_round_robin
	conn, err := grpc.Dial("etcd://Test/go_wallet_manage_svr",
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"weighted_round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	// 等待 1s 初始化完成
	time.Sleep(1 * time.Second)
	stub := pb.NewGoWalletManageSvrClient(conn)
	for i := 0; i < 7; i++ {
		rpcRsp, err := stub.ImportWallet(context.Background(), &pb.ImportWalletReq{
			PrivateKey: "0x12345",
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("rpcRsp: ", rpcRsp.GetAddress())
	}
}
