package collectors

import "node-exporter-lite/internal/metrics"

type Collector interface {
	Collect(getter metrics.MetricReader)
	RegisterMetrics(adder metrics.MetricAdder)
}
