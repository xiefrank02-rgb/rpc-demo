package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/xiefrank02-rgb/rpc-demo/hello" // 导入生成的 protobuf 代码

	"google.golang.org/grpc"
)

// 服务器实现 GreeterServer 接口
type server struct {
	hello.UnimplementedGreeterServer
}

// SayHello 方法实现
func (s *server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	return &hello.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func main() {
	// 创建监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册 Greeter 服务
	hello.RegisterGreeterServer(grpcServer, &server{})

	// 启动服务器
	fmt.Println("gRPC server listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
