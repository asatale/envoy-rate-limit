package main

import (
	"github.com/asatale/envoy-rate-limit/app/server/go/hello_world"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
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

	setLogLevel(appCfg.logLevel)

	lis, err := net.Listen("tcp", appCfg.addrValue)
	if err != nil {
		log.Fatal().Msgf("Failed to start listening: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			metricInterceptor,
			cancelInterceptor,
			delayInterceptor,
		),
	)
	hello_world.RegisterGreeterServer(s, &HelloWorldServer{})
	log.Info().Msgf("Server listening at %v", lis.Addr())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Info().Msgf("Received signal:", sig)
		s.GracefulStop()
		log.Info().Msgf("Server gracefully stopped")
	}()

	if err := startPrometheusServer(); err != nil {
		log.Fatal().Msgf("Failed to prometheus metric stub: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatal().Msgf("Failed to serve: %v", err)
	}
}
