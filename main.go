package main

import (
  "flag"
  "log"
  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "time"
)

var addr = flag.String("listen-address", ":9014", "The address to listen on for HTTP requests.")

var (
  vaultSealedGauge = prometheus.NewGauge(prometheus.GaugeOpts{
    Help: "The seal status of the vault instance",
    Namespace: "vault",
    Name: "sealed",
  })
)

func main() {
  flag.Parse()

  go func() {
    vaultSealedGauge.Set(1)
    time.Sleep(5 * time.Second)
  }()

  http.Handle("/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(*addr, nil))
}

func init() {
	prometheus.MustRegister(vaultSealedGauge)
}
