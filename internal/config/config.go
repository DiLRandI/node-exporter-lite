package config

import (
	"log/slog"

	"node-exporter-lite/internal/common"
)

const (
	defaultPort = 8080
)

type config struct {
	LogFilePath            *string
	LogLevel               *string
	Port                   *int
	PublishExporterMetrics *bool
}

func NewConfig() *config {
	return defaultConfig()
}

func defaultConfig() *config {
	return &config{
		LogFilePath:            common.MakePtr(""),
		LogLevel:               common.MakePtr(slog.LevelInfo.String()),
		Port:                   common.MakePtr(defaultPort),
		PublishExporterMetrics: common.MakePtr(false),
	}
}

func (c *config) ParseConfig() *config {
	pc := defaultConfig()

	if c.LogFilePath != nil && *c.LogFilePath != "" {
		pc.LogFilePath = c.LogFilePath
	}

	if c.LogLevel != nil && *c.LogFilePath != "" {
		pc.LogLevel = c.LogLevel
	}

	if c.Port != nil && *c.Port != 0 {
		pc.Port = c.Port
	}

	if c.PublishExporterMetrics != nil {
		pc.PublishExporterMetrics = c.PublishExporterMetrics
	}

	return pc
}
