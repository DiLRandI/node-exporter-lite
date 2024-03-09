package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var CPUTemperature = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current CPU temperature in Celsius",
})

func Register(register *prometheus.Registry) {
	register.MustRegister(CPUTemperature)
}
