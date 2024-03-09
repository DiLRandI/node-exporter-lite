package cpu

import (
	"context"
	"node-exporter-lite/internal/metrics"
	"os"
	"strconv"
	"strings"
	"time"
)

func Collect(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				readTemperature()
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func readTemperature() {
	tempData, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return
	}

	tempStr := strings.TrimSpace(string(tempData))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return
	}

	temp = temp / 1000

	metrics.CPUTemperature.Set(temp)
}
