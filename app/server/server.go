package main

import (
	"context"
	"log"
	"net"

	pb "github.com/envoy-rate-limit/app/proto/hello-world"
	"google.golang.org/grpc"
)


func main() {
	return 0
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}
