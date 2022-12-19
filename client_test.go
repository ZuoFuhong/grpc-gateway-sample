package main

import (
	"context"
	"fmt"
	pb "github.com/ZuoFuhong/grpc_gateway_best_practices/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func Test_Client(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:1025", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	stub := pb.NewGoEchoSvrClient(conn)
	rpcRsp, err := stub.Echo(context.Background(), &pb.EchoReq{
		Payload: "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rpcRsp.GetPayload())
}
