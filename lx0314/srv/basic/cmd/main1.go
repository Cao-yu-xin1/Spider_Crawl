package main

import (
	"context"
	"flag"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/nacos"
	__ "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/proto"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/basic/init1"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/handler"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	services, err := init1.GetServiceWithLoadBalancer(nacos.GlobalConfig.Consul.ServiceName)
	if err != nil {
		log.Printf("获取用户服务失败: %v", err)
	} else {
		log.Printf("获取到用户服务: %s, 地址: %s:%d", services.Service, services.Address, services.Port)
	}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &Server{})
	__.RegisterServiceServer(s, &handler.Server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")
	if err := init1.ConsulShutdown(); err != nil {
		log.Printf("Consul注销失败: %v", err)
	}
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.GracefulStop()
	log.Println("服务已关闭")
}
