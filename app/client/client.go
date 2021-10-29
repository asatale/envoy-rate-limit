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
}

func main() {
	addr := flag.String("addr", "localhost:50051", "Server address string")
	reqs := flag.Int("reqs", 100, "Number of request sent by a client")
	clients := flag.Int("clients", 1, "Number of clients")
	flag.Parse()

	var wg sync.WaitGroup
	for i := 0; i < *clients; i++ {
		wg.Add(1)
		go communicate(
			&clientConfig{
				server_addr:  *addr,
				client_id:    fmt.Sprintf("client-%d", i),
				num_requests: *reqs}, &wg)
	}
	wg.Wait()
}

func communicate(cfg *clientConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	latencies := make([]int64, cfg.num_requests)

	log.Printf("Client: %s Start", cfg.client_id)

	conn, err := grpc.Dial(cfg.server_addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	var wgSync sync.WaitGroup
	for i := 0; i < cfg.num_requests; i++ {
		wgSync.Add(1)
		go func(seqNum int, wgSync *sync.WaitGroup) {
			ctx, cancel := context.WithTimeout(context.Background(),
				time.Second*10)
			defer cancel()

			now := time.Now()
			start_msec := now.UnixMilli()
			r, err := c.SayHello(ctx, &pb.HelloRequest{ClientName: cfg.client_id,
				SeqNum: int32(seqNum)})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			if int32(seqNum) != r.GetSeqNum() {
				log.Fatalf("Invalid seqNum. Expected: %v, Received: %v",
					seqNum, r.GetSeqNum())
			}

			now = time.Now()
			end_msec := now.UnixMilli()
			latencies[seqNum] = end_msec - start_msec
			wgSync.Done()
		}(i, &wgSync)
	}
	wgSync.Wait()
	minRTT, maxRTT, sum := minMaxSum(latencies)
	log.Printf("Client: %s Done. Min RTT: %vms, Max RTT: %vms, Avg: %vms",
		cfg.client_id, minRTT, maxRTT, sum/int64(cfg.num_requests))
}

func minMaxSum(array []int64) (int64, int64, int64) {
	var sum int64
	max := array[0]
	min := array[0]

	for _, value := range array {
		sum = sum + value
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max, sum
}
