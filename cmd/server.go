package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	// 需要根据实际路径引入生成的hello包
	hello "github.com/xiefrank02-rgb/rpc-demo/api/hello/v1"
)

// 实现 Greeter 服务
type GreeterServer struct {
	hello.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	return &hello.HelloReply{
		Message: "Hello " + req.Name,
	}, nil
}

// 启动 gRPC 服务
func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	hello.RegisterGreeterServer(grpcServer, &GreeterServer{})

	// gRPC 服务反射
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	// 创建 Fiber 实例
	app := fiber.New()

	// 启动 gRPC 服务器
	go startGRPCServer()

	// Fiber HTTP 服务端口
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 设置 HTTP 请求处理器
	app.Post("/sayhello", func(c *fiber.Ctx) error {
		// 通过 HTTP 请求解析客户端传入的参数
		var req hello.HelloRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).SendString("Invalid Request")
		}

		// 连接到 gRPC 服务器
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) // 使用不安全的连接，仅用于开发
		if err != nil {
			return c.Status(500).SendString("Failed to connect to gRPC server")
		}
		defer conn.Close()

		// 创建 gRPC 客户端
		client := hello.NewGreeterClient(conn)

		// 调用 gRPC 方法
		resp, err := client.SayHello(c.Context(), &req)
		if err != nil {
			return c.Status(500).SendString("gRPC call failed")
		}

		// 返回响应
		return c.JSON(resp)
	})

	// 启动 HTTP 服务器
	log.Fatal(app.Listen(":8080"))
}
