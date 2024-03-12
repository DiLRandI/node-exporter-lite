package metrics

import (
	"regexp"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricAdder interface {
	AddGaugeMetric(subsystem, name, help string) error
}

type MetricReader interface {
	GetGaugeMetric(name string) (prometheus.Gauge, error)
}

type Register interface {
	Get() *prometheus.Registry
}

type MetricRegistry interface {
	MetricAdder
	MetricReader
	Register
}

type metricRegistry struct {
	metricsMap map[string]prometheus.Gauge
	register   *prometheus.Registry
	rwLock     *sync.RWMutex
	nameSpace  string
}

func NewRegistry(nameSpace string, publishExporterMetrics bool) MetricRegistry {
	var register *prometheus.Registry
	if publishExporterMetrics {
		r, ok := prometheus.DefaultRegisterer.(*prometheus.Registry)
		if !ok {
			panic("failed to get default registerer")
		}

		register = r
	} else {
		register = prometheus.NewRegistry()
	}

	return &metricRegistry{
		metricsMap: make(map[string]prometheus.Gauge),
		nameSpace:  nameSpace,
		rwLock:     new(sync.RWMutex),
		register:   register,
	}
}

func (m *metricRegistry) sanitizeMetricName(input string) string {
	// Replace invalid characters with underscores
	validName := regexp.MustCompile(`[^a-zA-Z0-9_:]`).ReplaceAllString(input, "_")
	// Ensure that the metric name starts with a letter or underscore

	return validName
}

func (m *metricRegistry) AddGaugeMetric(subsystem, name, help string) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	if m.metricsMap[name] != nil {
		return nil
	}

	gaugeMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      m.sanitizeMetricName(name),
		Help:      help,
		Namespace: m.nameSpace,
		Subsystem: subsystem,
	})

	if err := m.register.Register(gaugeMetric); err != nil {
		return MetricFailedToRegisterError(name)
	}

	m.metricsMap[name] = gaugeMetric

	return nil
}

func (m *metricRegistry) GetGaugeMetric(name string) (prometheus.Gauge, error) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()

	if m.metricsMap[name] == nil {
		return nil, MetricNotFoundError(name)
	}

	return m.metricsMap[name], nil
}

func (m *metricRegistry) Get() *prometheus.Registry {
	return m.register
}
