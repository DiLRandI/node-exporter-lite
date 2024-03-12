package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"node-exporter-lite/internal/common"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(port int, gather prometheus.Gatherer) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux(gather),
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			requestID := uuid.NewString()
			valueCtx := context.WithValue(ctx, common.RequestID{}, requestID)
			remoteAddr := c.RemoteAddr().String()
			valueCtx = context.WithValue(valueCtx, common.RemoteAddr{}, remoteAddr)

			return valueCtx
		},
	}

	return server
}

func mux(gather prometheus.Gatherer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.Handle("/metrics", promhttp.HandlerFor(gather, promhttp.HandlerOpts{}))

	// Register pprof handlers
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux
}
