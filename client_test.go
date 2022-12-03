package main

import (
	"context"
	"fmt"
	pb "github.com/ZuoFuhong/grpc_gateway_best_practices/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"testing"
)

func Test_Client(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:1024", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := pb.NewGoWalletManageSvrClient(conn)
	rpcRsp, err := client.ImportWallet(context.Background(), &pb.ImportWalletReq{
		PrivateKey: "0x12345",
	})
	if err != nil {
		// Resolve grpc errcode
		if rpcErr, ok := status.FromError(err); ok {
			fmt.Println(rpcErr)
		}
		return
	}
	fmt.Println("rpcRsp: ", rpcRsp.GetAddress())
}
