package main

import (
	"flag"
)

var (
	addrValue    *string
	delayValue   *int
	dprobValue   *int
	cancelOption *bool
	cprobValue   *int
)

func init() {
	addrValue = flag.String("addr", "0.0.0.0:50051", "Server address string")
	delayValue = flag.Int("delay", 20, "Response delay in millisecond")
	dprobValue = flag.Int("dprob", 20, "Delay Probability")
	cancelOption = flag.Bool("cancel", false, "Cancel RPC with cancel-probability")
	cprobValue = flag.Int("cprob", 20, "Cancel Probability")
	flag.Parse()
}
