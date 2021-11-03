package main

import (
	"context"
	"flag"
	pb "github.com/asatale/envoy-rate-limit/app/proto/hello_world"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"time"
)

type server struct {
	pb.UnimplementedGreeterServer
}

var (
	addr      *string
	delay     *int
	variance  *int
	randomGen rand.Source
)

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	if *delay != 0 {
		r1 := rand.New(randomGen)
		var v int
		if *variance != 0 {
			v = r1.Intn(*variance)
		}

		rsp_delay := time.Duration(*delay)

		if v%2 == 0 {
			rsp_delay = rsp_delay + time.Duration(v)
		} else {
			rsp_delay = rsp_delay - time.Duration(v)
		}
		time.Sleep(rsp_delay * time.Millisecond)
	}
	return &pb.HelloReply{ClientName: "Hello " + in.GetClientName(), SeqNum: in.GetSeqNum()}, nil
}

func main() {
	addr = flag.String("addr", "0.0.0.0:50051", "Server address string")
	delay = flag.Int("rsp_delay", 10, "Response delay in millisecond")
	variance = flag.Int("variance", 2, "Response time randomized variance in millisecond")

	flag.Parse()

	// Initialize random number generator
	randomGen = rand.NewSource(time.Now().UnixNano())

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
