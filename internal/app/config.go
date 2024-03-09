package app

import (
	"log/slog"
	"node-exporter-lite/internal/common"
)

const defaultPort = 8080
const defaultLogFilePath = "/var/log/node-exporter-lite.log"

type config struct {
	LogFilePath            *string
	LogLevel               *string
	Port                   *int
	PublishExporterMetrics *bool
}

func DefaultConfig() *config {
	return &config{
		LogFilePath:            common.MakePtr(defaultLogFilePath),
		LogLevel:               common.MakePtr(slog.LevelInfo.String()),
		Port:                   common.MakePtr(defaultPort),
		PublishExporterMetrics: common.MakePtr(false),
	}
}

func (c *config) ParseConfig() *config {
	pc := DefaultConfig()

	if c.LogFilePath != nil && *c.LogFilePath != "" {
		pc.LogFilePath = c.LogFilePath
	}

	if c.LogLevel != nil && *c.LogFilePath != "" {
		pc.LogLevel = c.LogLevel
	}

	return pc
}
