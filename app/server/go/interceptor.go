package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"time"
)

func cancelInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if *cancelOption && *cprobValue > 0 {
		randNumber := rand.Intn(101)

		if randNumber <= *cprobValue {
			log.Printf("Cancelling RPC")
			_, cancel := context.WithCancel(ctx)
			defer cancel()
			return nil, status.Errorf(codes.ResourceExhausted, "%s is cancelled by middleware", info.FullMethod)
		}
	}
	return handler(ctx, req)
}

func delayInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if *delayValue > 0 && *dprobValue > 0 {
		randNumber := rand.Intn(101)
		if randNumber <= *dprobValue {
			log.Printf("Delayed RPC response")
			time.Sleep(time.Duration(*delayValue))
		}
	}
	return handler(ctx, req)
}
