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

const (
	STARTUP_TIMEOUT = 1 // 1 sec wait on startup before starting traffic
	RPC_TIMEOUT     = 5 // 5 sec connection timeout
)

type clientConfig struct {
	thread_id    string
	server_addr  string
	num_requests int
}

func main() {
	addr := flag.String("addr", "0.0.0.0:50051", "Server address string")
	reqs := flag.Int("reqs", 100, "Number of request sent by a client")
	num_threads := flag.Int("num_threads", 1, "Number of concurrent threads")
	flag.Parse()

	log.Printf("Waiting for servers to come up")
	time.Sleep(STARTUP_TIMEOUT * time.Second)

	log.Printf("Starting %v threads", *num_threads)
	var wg sync.WaitGroup
	for i := 0; i < *num_threads; i++ {
		wg.Add(1)
		go communicate(
			&clientConfig{
				server_addr:  *addr,
				thread_id:    fmt.Sprintf("thread-%d", i),
				num_requests: *reqs}, &wg)
	}
	wg.Wait()
}

func communicate(cfg *clientConfig, clientWg *sync.WaitGroup) {
	defer clientWg.Done()
	var (
		conn       *grpc.ClientConn
		errCnt     int64
		errCntLock sync.Mutex
		msgWg      sync.WaitGroup
	)

	latencies := make([]int64, cfg.num_requests)
	log.Printf("Thread: %s Started", cfg.thread_id)

	conn, err := grpc.Dial(cfg.server_addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Thread: %s could not connect: %v", cfg.thread_id, err)
	}
	defer conn.Close()

	clnt := pb.NewGreeterClient(conn)

	for i := 0; i < cfg.num_requests; i++ {
		msgWg.Add(1)

		go func(seqNum int, wgSync *sync.WaitGroup, errCntLock *sync.Mutex) {
			defer wgSync.Done()

			ctx, cancel := context.WithTimeout(context.Background(),
				time.Second*RPC_TIMEOUT)
			defer cancel()

			now := time.Now()
			start_msec := now.UnixMilli()

			_, err := clnt.SayHello(ctx, &pb.HelloRequest{ClientName: cfg.thread_id,
				SeqNum: int32(seqNum)})

			if err != nil {
				log.Printf("RPC Error")
				errCntLock.Lock()
				errCnt = errCnt + 1
				errCntLock.Unlock()
				return
			}

			now = time.Now()
			end_msec := now.UnixMilli()
			latencies[seqNum] = end_msec - start_msec

		}(i, &msgWg, &errCntLock)
	}
	msgWg.Wait()
	minRTT, maxRTT, sum := minMaxSum(latencies)
	log.Printf("Thread: %s Done. Min RTT: %vms, Max RTT: %vms, Avg: %vms, ErrRsp: %v",
		cfg.thread_id, minRTT, maxRTT, sum/int64(cfg.num_requests), errCnt)
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
