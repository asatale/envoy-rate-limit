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

func metricInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	total_rpc_metric.Inc()
	return handler(ctx, req)
}

func cancelInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if appCfg.cancelOption && appCfg.cprobValue > 0 {
		randNumber := rand.Intn(101)

		if randNumber <= appCfg.cprobValue {
			log.Printf("Cancelling RPC")
			cancel_rpc_metric.Inc()
			return nil, status.Errorf(codes.ResourceExhausted, "%s is cancelled by middleware", info.FullMethod)
		}
	}
	return handler(ctx, req)
}

func delayInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if appCfg.delayValue > 0 && appCfg.dprobValue > 0 {
		randNumber := rand.Intn(101)
		if randNumber <= appCfg.dprobValue {
			log.Printf("Delayed RPC response")
			delayed_rpc_metric.Inc()
			time.Sleep(time.Duration(appCfg.delayValue) * time.Millisecond)
		}
	}
	return handler(ctx, req)
}
