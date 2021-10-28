package main

import (
	"context"
	"log"
	"net"

	pb "github.com/asatale/envoy-rate-limit/proto/hello_world"
	"google.golang.org/grpc"
)


func main() {
	return 0
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}
