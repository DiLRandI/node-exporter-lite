package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"node-exporter-lite/internal/collectors/thermal"
	"node-exporter-lite/internal/config"
	"node-exporter-lite/internal/metrics"
	"node-exporter-lite/internal/server"
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

	metricRegistry := metrics.NewRegistry("node_exporter_lite", *config.PublishExporterMetrics)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	thermal := thermal.New(logger, metricRegistry, metricRegistry, time.Second*5)
	thermal.Collect(ctx)

	server := server.NewServer(*config.Port, metricRegistry.Get())

	server.RegisterOnShutdown(func() {
		logger.WarnContext(ctx, "server is shutting down")
	})

	logger.InfoContext(ctx, "server is starting", "port", *config.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.ErrorContext(ctx, "server failed to start", "error", err)
	}

	logger.Warn("application is exiting")
}
