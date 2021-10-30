package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/asatale/envoy-rate-limit/app/proto/hello_world"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	STARTUP_TIMEOUT = 1 // 1 sec wait on startup before starting traffic
	RPC_TIMEOUT     = 5 // 5 sec connection timeout
)

type clientConfig struct {
	thread_id      string
	server_addr    string
	num_requests   int
	burst_size     int
	burst_interval int
}

func main() {
	server_addr := flag.String("server_addr",
		"localhost:50051",
		"Server address string.")

	num_threads := flag.Int("num_threads",
		10,
		"Number of concurrent threads.")

	reqs_per_thread := flag.Int("reqs_per_thread",
		100,
		"Total Number of request sent by a client.")

	burst_size := flag.Int("burst_size",
		100, "Number of RPCs in single burst interval.")

	burst_interval := flag.Int("burst_interval",
		100,
		"Burst internval in millisecond")

	flag.Parse()

	log.Printf("Waiting for servers to come up")
	time.Sleep(STARTUP_TIMEOUT * time.Second)

	var wg sync.WaitGroup
	for i := 0; i < *num_threads; i++ {
		wg.Add(1)
		go communicate(
			&clientConfig{
				server_addr:    *server_addr,
				thread_id:      fmt.Sprintf("thread-%d", i),
				num_requests:   *reqs_per_thread,
				burst_size:     *burst_size,
				burst_interval: *burst_interval,
			}, &wg)
	}
	wg.Wait()
}

func randomize_start() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	delay := time.Duration(r1.Intn(100))
	time.Sleep(delay * time.Millisecond)
}

func communicate(cfg *clientConfig, clientWg *sync.WaitGroup) {
	defer clientWg.Done()
	var (
		conn       *grpc.ClientConn
		errCnt     int64
		successCnt int64
		cntLock    sync.Mutex
		msgWg      sync.WaitGroup
	)

	randomize_start()

	latencies := make([]int64, cfg.num_requests)
	log.Printf("Thread: %s Started", cfg.thread_id)

	conn, err := grpc.Dial(cfg.server_addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Thread: %s could not connect: %v", cfg.thread_id, err)
	}
	defer conn.Close()

	clnt := pb.NewGreeterClient(conn)

	limiter := time.Tick(time.Duration(cfg.burst_interval) * time.Millisecond)

	for i, b := 0, 0; i < cfg.num_requests; i++ {
		msgWg.Add(1)
		go func(seqNum int, wgSync *sync.WaitGroup, cntLock *sync.Mutex) {
			defer wgSync.Done()

			ctx, cancel := context.WithTimeout(context.Background(),
				time.Second*RPC_TIMEOUT)
			defer cancel()

			now := time.Now()
			start_msec := now.UnixMilli()

			_, err := clnt.SayHello(ctx, &pb.HelloRequest{ClientName: cfg.thread_id,
				SeqNum: int32(seqNum)})

			if err != nil {
				log.Printf("RPC Error: %v", err)
				cntLock.Lock()
				errCnt = errCnt + 1
				cntLock.Unlock()
				return
			} else {
				cntLock.Lock()
				successCnt = successCnt + 1
				cntLock.Unlock()
			}

			now = time.Now()
			end_msec := now.UnixMilli()
			latencies[seqNum] = end_msec - start_msec
			if latencies[seqNum] < 0 { // Time skew
				latencies[seqNum] = 0
			}
		}(i, &msgWg, &cntLock)

		if b == cfg.burst_size {
			<-limiter
			b = 0
		} else {
			b = b + 1
		}
	}
	msgWg.Wait()
	minRTT, maxRTT, sum := minMaxSum(latencies)
	log.Printf("Thread: %s Done. Min RTT: %vms, Max RTT: %vms, Avg: %vms, SucessRsp: %v ErrRsp: %v",
		cfg.thread_id, minRTT, maxRTT, sum/int64(cfg.num_requests), successCnt, errCnt)
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
