package main

import (
	"github.com/asatale/envoy-rate-limit/app/server/go/hello_world"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"time"
)

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

func main() {

	lis, err := net.Listen("tcp", *addrValue)
	if err != nil {
		log.Fatalf("Failed to start listening: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			cancelInterceptor,
			delayInterceptor,
		),
	)
	hello_world.RegisterGreeterServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
