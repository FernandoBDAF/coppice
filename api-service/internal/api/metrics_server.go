package api

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fernandobarroso/microservices/api-service/internal/config"
)

// NewMetricsServer builds an HTTP server that exposes Prometheus metrics on
// its own port (cfg.Metrics.Port, default 8081), separate from the main API
// port. Callers are responsible for starting it (ListenAndServe in a
// goroutine) and shutting it down alongside the main server.
func NewMetricsServer(cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Metrics.Port),
		Handler: mux,
	}
}
