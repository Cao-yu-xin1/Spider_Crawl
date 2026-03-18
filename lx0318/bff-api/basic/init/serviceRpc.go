package init

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"lx0318/config"
	__ "lx0318/proto"
)

func init() {
	ServiceRpc()
}

func ServiceRpc() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	//c := pb.NewGreeterClient(conn)
	config.ServiceClient = __.NewServiceClient(conn)
}
