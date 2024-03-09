package collector

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWrite(t *testing.T) {
	c := NewCollector(slog.Default())
	c.Write("test", "test")

	data := c.Read()
	assert.Equal(t, data, map[string]string{"test": "test"})

}
