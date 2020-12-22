package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/safchain/ethtool"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	ethHandle, err := ethtool.NewEthtool()

	// Create a map of the gauges for every single metric that we can periodically update
	// We can filter which metrics get scraped on the scraper side so we should expose all of them
	metricsMap := make(map[string]prometheus.Gauge)
	stats, err := ethHandle.Stats("eth0")
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	for metricName := range stats {
		log.Printf("%#v\n", metricName)
		// TODO: Make labels for these configurable
		// TODO: Make namespace configurable
		metricsMap[metricName] = promauto.NewGauge(prometheus.GaugeOpts{
			Name: metricName,
		})
	}

	// We don't have an exit case for this goroutine because we want it to run
	// forever anyways
	go func() {
		for {
			if err != nil {
				log.Fatalf("err: %v\n", err)
			}
			defer ethHandle.Close()

			stats, err := ethHandle.Stats("eth0")
			if err != nil {
				log.Fatalf("err: %v\n", err)
			}

			for metricName, metricValue := range stats {
				metricsMap[metricName].Set(float64(metricValue))
			}

			// TODO: Make delay configurable
			time.Sleep(5 * time.Second)
		}
	}()

	// No graceful shutdown for this web server because we don't really need it
	// to
	srv := &http.Server{
		// TODO: Make address configurable
		Addr:    ":8081",
		Handler: mux,
	}
	srv.ListenAndServe()
}
