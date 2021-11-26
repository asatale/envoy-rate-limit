package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

var (
	total_rpc_metric, cancel_rpc_metric, delayed_rpc_metric prometheus.Counter
)

func init() {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	total_rpc_metric = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "go_grpc_server_total_requests",
		Help:        "Total number of RPC requests received",
		ConstLabels: prometheus.Labels{"host": hostName, "app": "Go GRPCserver"},
	})
	cancel_rpc_metric = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "go_grpc_server_cancelled_requests",
		Help:        "Number of cancelled RPCs",
		ConstLabels: prometheus.Labels{"host": hostName, "app": "Go GRPCserver"},
	})
	delayed_rpc_metric = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "go_grpc_server_delayed_requests",
		Help:        "Number of delated RPC responses",
		ConstLabels: prometheus.Labels{"host": hostName, "app": "Go GRPCserver"},
	})
}

func startPrometheusServer() error {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(total_rpc_metric)
	prometheus.MustRegister(cancel_rpc_metric)
	prometheus.MustRegister(delayed_rpc_metric)
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatalf("Failed to start prometheus endpoint")
		}
	}()
	return nil
}
