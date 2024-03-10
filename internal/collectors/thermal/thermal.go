package thermal

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"node-exporter-lite/internal/collectors"
	"node-exporter-lite/internal/metrics"
)

const (
	thermalZonePath = "/sys/class/thermal/"
	baseMetricPath  = "thermal"
)

type ThermalCollector struct {
	metricAdder  metrics.MetricAdder
	metricGetter metrics.MetricReader
	thermalZones map[string]string
	logger       *slog.Logger
	interval     time.Duration
}

func New(
	logger *slog.Logger,
	metricAdder metrics.MetricAdder,
	metricGetter metrics.MetricReader,
	interval time.Duration,
) collectors.Collector {
	thermalCollector := &ThermalCollector{
		thermalZones: make(map[string]string),
		logger:       logger,
		interval:     interval,
		metricAdder:  metricAdder,
		metricGetter: metricGetter,
	}

	thermalCollector.readAllThermalZones()

	return thermalCollector
}

func (c *ThermalCollector) readAllThermalZones() {
	dir, err := os.ReadDir(thermalZonePath)
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {
		if !strings.HasPrefix(entry.Name(), "thermal_zone") {
			continue
		}

		temperaturePath := thermalZonePath + entry.Name() + "/temp"
		typePath := thermalZonePath + entry.Name() + "/type"

		typeContent, err := c.readContent(typePath)
		if err != nil {
			panic(err)
		}

		metricName := c.buildMetricName(entry.Name(), typeContent)
		c.thermalZones[metricName] = temperaturePath
	}

	c.registerMetric()
}

func (c *ThermalCollector) registerMetric() {
	for typ := range c.thermalZones {
		if err := c.metricAdder.AddGaugeMetric(baseMetricPath, typ, "Temperature in celsius"); err != nil {
			c.logger.Error("failed to add thermal metric", "error", err)
		}
	}
}

func (c *ThermalCollector) readContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func (c *ThermalCollector) buildMetricName(name, typ string) string {
	return name + "_" + strings.ReplaceAll(strings.ToLower(typ), " ", "_")
}

func (c *ThermalCollector) Collect(ctx context.Context) {
	go c.collect(ctx)
}

func (c *ThermalCollector) collect(ctx context.Context) {
	c.logger.Info("thermal collector started")
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("thermal collector stopped")
			return
		case <-time.After(c.interval):
			c.collectThermalZones()
		}
	}
}

func (c *ThermalCollector) collectThermalZones() {
	for typ, path := range c.thermalZones {
		metric, err := c.metricGetter.GetGaugeMetric(typ)
		if err != nil {
			c.logger.Error("failed to get thermal metric", "name", typ, "error", err)
			continue
		}

		tempStr, err := c.readContent(path)
		if err != nil {
			c.logger.Error("failed to read thermal zone", "name", typ, "path", path, "error", err)
			continue

		}

		temperature, err := strconv.Atoi(tempStr)
		if err != nil {
			c.logger.Error("failed to convert temperature", "name", typ, "temperature", tempStr, "error", err)
			continue
		}

		temp := float64(temperature) / 1000.0

		metric.Set(temp)
	}
}
