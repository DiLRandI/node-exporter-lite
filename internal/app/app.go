package app

import (
	"io"
	"log/slog"
	"net/http"
	"node-exporter-lite/internal/common"
	"node-exporter-lite/internal/service/collector"
)

type App struct {
	logger          *slog.Logger
	metricCollector collector.MetricReader
}

func NewApp(logger *slog.Logger, metricCollector collector.MetricReader) *App {
	app := &App{
		logger:          logger,
		metricCollector: metricCollector,
	}

	return app
}

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("use /metrics"))
		if err != nil {
			a.logger.Error("error writing response",
				"requestID", r.Context().Value(common.RequestID{}),
				"remoteAddr", r.Context().Value(common.RemoteAddr{}),
				"error", err)
		}
	}))
	mux.Handle("GET /metrics", a.metrics())

	return mux
}

func (a *App) metrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		a.logger.Info("/metrics",
			"requestID", r.Context().Value(common.RequestID{}),
			"remoteAddr", r.Context().Value(common.RemoteAddr{}))
		a.Write(w)
	}
}

func (a *App) Write(w io.Writer) {
	data := a.metricCollector.Read()
	for k, v := range data {
		_, err := w.Write([]byte(k + " " + v + "\n"))
		if err != nil {
			a.logger.Error("error writing response", "error", err)
		}
	}
}
