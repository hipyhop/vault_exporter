package main

import (
  "flag"
  "log"
  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "time"
)

var addr = flag.String("listen-address", ":9410", "The address to listen on for HTTP requests.")
var checkInterval = flag.Int("check-interval", 20, "How frequently, in seconds, to check vault metrics.")

var (
  vaultSealedGauge = prometheus.NewGauge(prometheus.GaugeOpts{
    Help: "The seal status of the vault instance",
    Namespace: "vault",
    Name: "sealed",
  })
)

func main() {
  flag.Parse()

  // Fake vault status
  go func() {
    sleepDuration := time.Duration(*checkInterval) * time.Second
    for {
      vaultSealedGauge.Set(1)
      time.Sleep(sleepDuration)
    }
  }()

  http.Handle("/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(*addr, nil))
}

func init() {
	prometheus.MustRegister(vaultSealedGauge)
}
