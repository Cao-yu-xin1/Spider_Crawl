package main

import (
	"flag"
	"log"
	__ "lx0318/proto"
	_ "lx0318/rpc/basic/init"
	"lx0318/rpc/handler"
	"net"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &handler.Server{})
	__.RegisterServiceServer(s, &handler.Server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
