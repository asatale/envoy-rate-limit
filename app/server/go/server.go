package main

import (
	"github.com/asatale/envoy-rate-limit/app/server/go/hello_world"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

func main() {

	lis, err := net.Listen("tcp", appCfg.addrValue)
	if err != nil {
		log.Fatalf("Failed to start listening: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			metricInterceptor,
			cancelInterceptor,
			delayInterceptor,
		),
	)
	hello_world.RegisterGreeterServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Received signal:", sig)
		s.GracefulStop()
		log.Println("Server gracefully stopped")
	}()

	if err := startPrometheusServer(); err != nil {
		log.Fatalf("Failed to prometheus metric stub: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
