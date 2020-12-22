package main

import (
	"flag"
	"log"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/safchain/ethtool"
)

func main() {
	statsHost := flag.String("host", "", "the datadog agent host and port")
	flag.Parse()

	statsd, err := statsd.New(*statsHost, statsd.WithNamespace("aws.ec2.ethtool"))
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	for {
		ethHandle, err := ethtool.NewEthtool()
		if err != nil {
			log.Fatalf("err: %v\n", err)
		}
		defer ethHandle.Close()

		stats, err := ethHandle.Stats("eth0")
		if err != nil {
			log.Fatalf("err: %v\n", err)
		}

		track := []string{"bw_in_allowance_exceeded", "bw_out_allowance_exceeded", "pps_allowance_exceeded", "conntrack_allowance_exceeded", "linklocal_allowance_exceeded"}

		for _, stat := range track {
			err = statsd.Count(stat, int64(stats[stat]), nil, 1)
			if err != nil {
				log.Printf("err posting gauge %v: %v\n", stat, err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
