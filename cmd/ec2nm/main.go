package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/safchain/ethtool"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Caller().Logger()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	ethHandle, err := ethtool.NewEthtool()

	// Create a map of the gauges for every single metric that we can periodically update
	// We can filter which metrics get scraped on the scraper side so we should expose all of them
	metricsMap := make(map[string]prometheus.Gauge)
	stats, err := ethHandle.Stats("eth0")
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't get eth0 stats")
	}

	log.Info().Msg("Initializing metrics map")
	for metricName := range stats {
		// TODO: Make labels for these configurable
		// TODO: Make namespace configurable
		metricsMap[metricName] = promauto.NewGauge(prometheus.GaugeOpts{
			Name: metricName,
		})
	}

	// We don't have an exit case for this goroutine because we want it to run
	// forever anyways
	go func() {
		log.Info().Msg("Starting metrics collection loop")
		for {
			stats, err := ethHandle.Stats("eth0")
			if err != nil {
				log.Fatal().Err(err).Msg("Couldn't get eth0 stats")
			}

			for metricName, metricValue := range stats {
				metricsMap[metricName].Set(float64(metricValue))
			}

			// TODO: Make delay configurable
			time.Sleep(5 * time.Second)
		}
	}()

	log.Info().Msg("Starting http server to serve prometheus metrics")
	// No graceful shutdown for this web server because we don't really need
	srv := &http.Server{
		// TODO: Make address configurable
		Addr:    ":8081",
		Handler: mux,
	}
	srv.ListenAndServe()
	log.Info().Msg("Application is exiting")
}
