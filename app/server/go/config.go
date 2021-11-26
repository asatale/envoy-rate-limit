package main

import (
	"flag"
)

type appConfig struct {
	addrValue    string
	delayValue   int
	dprobValue   int
	cancelOption bool
	cprobValue   int
}

var appCfg appConfig

func init() {
	addrValue := flag.String("addr", "0.0.0.0:50051", "Server address string")
	delayValue := flag.Int("delay", 20, "Response delay in millisecond")
	dprobValue := flag.Int("dprob", 20, "Delay Probability")
	cancelOption := flag.Bool("cancel", false, "Cancel RPC with cancel-probability")
	cprobValue := flag.Int("cprob", 20, "Cancel Probability")
	flag.Parse()

	appCfg.addrValue = *addrValue
	appCfg.delayValue = *delayValue
	appCfg.dprobValue = *dprobValue
	appCfg.cancelOption = *cancelOption
	appCfg.cprobValue = *cprobValue
}
