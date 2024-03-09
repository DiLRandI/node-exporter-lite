package collector

import (
	"context"
	"log/slog"
	"maps"
	"sync"
)

type MetricReader interface {
	Read() map[string]string
}

type MetricWriter interface {
	Write(key, value string)
}

type Collector interface {
	Start(ctx context.Context)
}

type collector struct {
	inMemory map[string]string
	rwMutex  *sync.RWMutex

	logger *slog.Logger
}

func NewCollector(logger *slog.Logger) Collector {
	return &collector{
		inMemory: make(map[string]string),
		rwMutex:  new(sync.RWMutex),
		logger:   logger,
	}
}

func (c *collector) Read() map[string]string {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	c.logger.Debug("reading in-memory data")

	return maps.Clone(c.inMemory)
}

func (c *collector) Write(key, value string) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.logger.Debug("writing in-memory data", "key", key, "value", value)

	c.inMemory[key] = value
}

func (c *collector) Start(ctx context.Context) {
	c.logger.Info("starting collectors")
	defer c.logger.Info("stopping collectors")
}
