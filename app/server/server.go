package main

import (
	"context"
	"flag"
	pb "github.com/asatale/envoy-rate-limit/app/proto/hello_world"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{ClientName: "Hello " + in.GetClientName(), SeqNum: in.GetSeqNum()}, nil
}

func main() {
	addr := flag.String("addr", "0.0.0.0:50051", "Server address string")
	flag.Parse()

	log.Printf("Server listening on %v", *addr)
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to start listening: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
