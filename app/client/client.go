package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/asatale/envoy-rate-limit/app/proto/hello_world"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type clientConfig struct {
	client_id    string
	server_addr  string
	num_requests int
	burstSize    int
}

func main() {
	addr := flag.String("addr", "localhost:50051", "Server address string")
	reqs := flag.Int("reqs", 100, "Number of request sent by a client")
	clients := flag.Int("clients", 1, "Number of clients")
	burstSize := flag.Int("burst-size", 10, "Request Burst Size")
	flag.Parse()

	var wg sync.WaitGroup
	for i := 0; i < *clients; i++ {
		wg.Add(1)
		go communicate(
			&clientConfig{
				server_addr:  *addr,
				client_id:    fmt.Sprintf("client-%d", i),
				num_requests: *reqs,
				burstSize:    *burstSize}, &wg)
	}
	wg.Wait()
}

func communicate(cfg *clientConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Client: %s Start", cfg.client_id)

	conn, err := grpc.Dial(cfg.server_addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	for i := 1; i <= cfg.num_requests; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{ClientName: cfg.client_id,
			SeqNum: int32(i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Received: %s, SeqNum: %d", r.GetClientName(), r.GetSeqNum())
	}
	log.Printf("Client: %s Done", cfg.client_id)
}
