package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "time"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  vaultApi "github.com/hashicorp/vault/api"
)

var addr = flag.String("listen-address", ":9410", "The address to listen on for HTTP requests.")
var checkInterval = flag.Int("check-interval", 20, "How frequently, in seconds, to check vault metrics.")

var (
  vaultSealedGauge = prometheus.NewGauge(prometheus.GaugeOpts{
    Help: "The seal status of the vault instance",
    Namespace: "vault",
    Name: "sealed",
  })
  vaultQueryErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
    Help: "The number of failed api requests to the vault instance",
    Namespace: "vault",
    Name: "query_error_total",
  })
)

func main() {
  flag.Parse()

  vault, err := initVaultClient()
  if err != nil {
    log.Fatal(err)
  }

  go func() {
    sleepDuration := time.Duration(*checkInterval) * time.Second
    for {
      collectMetrics(vault)
      time.Sleep(sleepDuration)
    }
  }()

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `<h1>vault-exporter</h1><a href="/metrics">/metrics</a>`)
  })
  http.Handle("/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(*addr, nil))
}

func collectMetrics(v *vaultApi.Client) {
  health, err := v.Sys().Health()
  if err != nil {
    vaultQueryErrorCounter.Inc()
    fmt.Printf("Failed to check vault health: %s\n", err)
    return
  }

  vaultSealedGauge.Set(bool2Float64(health.Sealed))
}

func initVaultClient() (*vaultApi.Client, error) {
  vaultConfig := vaultApi.DefaultConfig()
  // Configure with ENV vars

  return vaultApi.NewClient(vaultConfig)
}

func bool2Float64(b bool)(float64) {
  if b {
    return 1
  }
  return 0
}

func init() {
  prometheus.MustRegister(vaultSealedGauge)
  prometheus.MustRegister(vaultQueryErrorCounter)
}
