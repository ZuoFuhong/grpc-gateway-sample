package main

import (
	"context"
	"github.com/ZuoFuhong/grpc-naming-monica/registry"
	pb "github.com/ZuoFuhong/grpc_gateway_best_practices/proto"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:1025")
	if err != nil {
		log.Fatal(err)
	}
	serviceImpl := new(GoWalletManageSvrImpl)
	s := grpc.NewServer()
	pb.RegisterGoWalletManageSvrServer(s, serviceImpl)
	// 1025 端口启动 gRPC Server
	log.Println("Serving gRPC on 127.0.0.1:1025")
	go func() {
		// 服务注册
		if err := registry.NewRegistry(&registry.Config{
			Token:       "18ee7064-3cdd-4ed5-a139-fd8d9add5847",
			Namespace:   "Test",
			ServiceName: "go_wallet_manage_svr",
			IP:          "127.0.0.1",
			Port:        1025,
			Weight:      100,
			Metadata:    "[]",
			HeartBeat:   5,
		}).Register(); err != nil {
			log.Fatalln(err)
		}
		log.Fatalln(s.Serve(lis))
	}()

	// 创建一个连接到 gRPC 服务器的客户端连接
	// gRPC-Gateway 就是通过它来代理请求（将 HTTP 请求转为 RPC 请求）
	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1:1025",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gw := runtime.NewServeMux()
	err = pb.RegisterGoWalletManageSvrHandler(context.Background(), gw, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    "127.0.0.1:8090",
		Handler: gw,
	}
	// 8090 端口提供 gRPC-Gateway 服务
	log.Println("Serving gRPC-Gateway on http://127.0.0.1:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

type GoWalletManageSvrImpl struct {
}

// CreateWallet 创建钱包
func (s *GoWalletManageSvrImpl) CreateWallet(ctx context.Context, _ *pb.CreateWalletReq) (*pb.CreateWalletRsp, error) {
	address := uuid.New().String()
	rsp := &pb.CreateWalletRsp{
		Address: address,
	}
	return rsp, nil
}

// ImportWallet 导入钱包
func (s *GoWalletManageSvrImpl) ImportWallet(ctx context.Context, req *pb.ImportWalletReq) (*pb.ImportWalletRsp, error) {
	log.Println("privKey: ", req.GetPrivateKey())
	address := uuid.New().String()
	rsp := &pb.ImportWalletRsp{
		Address: address,
	}
	return rsp, nil
}
