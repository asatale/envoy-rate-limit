package main

import (
	"context"
	"github.com/asatale/envoy-rate-limit/app/server/go/hello_world"
)

type server struct {
	hello_world.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *hello_world.HelloRequest) (*hello_world.HelloReply, error) {
	return &hello_world.HelloReply{ClientName: "Hello " + in.GetClientName(), SeqNum: in.GetSeqNum()}, nil
}
