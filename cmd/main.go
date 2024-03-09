package main

import (
	"context"
	"flag"
	"log/slog"
	"node-exporter-lite/internal/config"
	"node-exporter-lite/internal/metrics"
	"node-exporter-lite/internal/server"
	"node-exporter-lite/internal/service/collector"
	"node-exporter-lite/internal/service/collectors/cpu"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	config := config.NewConfig()

	// Parse the command line flags
	config.LogFilePath = flag.String("log-path", *config.LogFilePath, "Log file path")
	config.LogLevel = flag.String("level", *config.LogLevel, "Log level")
	config.Port = flag.Int("port", *config.Port, "Port to listen on")
	config.PublishExporterMetrics = flag.Bool(
		"publish-exporter-metrics", *config.PublishExporterMetrics, "Publish exporter metrics")

	flag.Parse()
	config = config.ParseConfig()
	level := new(slog.Level)
	if err := level.UnmarshalText([]byte(*config.LogLevel)); err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: *level,
	}))
	logger.Info("starting application",
		"port", *config.Port,
		"logFilePath", *config.LogFilePath,
		"logLevel", *config.LogLevel,
		"publishExporterMetrics", *config.PublishExporterMetrics)

	var register *prometheus.Registry
	if *config.PublishExporterMetrics {
		r, ok := prometheus.DefaultRegisterer.(*prometheus.Registry)
		if !ok {
			panic("failed to get default registerer")
		}

		register = r
	} else {
		register = prometheus.NewRegistry()
	}

	metrics.Register(register)

	collector := collector.NewCollector(logger)
	collector.Start(context.Background())

	server := server.NewServer(*config.Port, register)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server.RegisterOnShutdown(func() {
		logger.WarnContext(ctx, "server is shutting down")
	})

	cpu.Collect(ctx)

	logger.InfoContext(ctx, "server is starting", "port", *config.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.ErrorContext(ctx, "server failed to start", "error", err)
	}

	logger.Warn("application is exiting")
}
