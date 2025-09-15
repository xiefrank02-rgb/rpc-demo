package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw "github.com/xiefrank02-rgb/rpc-demo/api/hello/v1"
	hello "github.com/xiefrank02-rgb/rpc-demo/api/hello/v1"
	// 正确导入 grpc-gateway 的 runtime 包
)

// gRPC 服务实现
type GreeterServer struct {
	hello.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	return &hello.HelloReply{Message: "Hello " + req.Name}, nil
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	hello.RegisterGreeterServer(grpcServer, &GreeterServer{})
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startGateway() {
	ctx := context.Background()

	// 使用 grpc-gateway 提供的 runtime.ServeMux，而不是 http.ServeMux
	mux := runtime.NewServeMux() // 使用 grpc-gateway 的 ServeMux

	// 设置 grpc 连接选项
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 注册 gRPC 服务的 HTTP Gateway
	if err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// 启动 HTTP Gateway
	log.Println("HTTP Gateway listening on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}
}

func main() {
	go startGRPCServer()
	go startGateway()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 修正 swagger 路径
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/doc.json", // 使用相对路径
	}))

	// 提供 swagger.json
	app.Static("/swagger/doc.json", "./api/hello/v1/hello.swagger.json")

	log.Fatal(app.Listen(":8080"))
}
