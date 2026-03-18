package init

import (
	"flag"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/bff-api/basic/config"
	__ "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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
