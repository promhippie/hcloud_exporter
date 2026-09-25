package exporter

import (
	"testing"

	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Lists parsed from a flag or an environment variable carry spare capacity, so
// appending the disk and interface labels in place would give both metric
// families whichever label got appended last.
func TestServerMetricsCollectorKeepsDerivedLabelsApart(t *testing.T) {
	collector := NewServerMetricsCollector(
		testLogger(),
		nil,
		nil,
		nil,
		config.Target{
			ServerMetrics: config.ServerMetrics{Labels: spareCapacity(config.ServerMetricsLabels())},
		},
	)

	assert.Contains(t, collector.CPU.String(), "variableLabels: {id,name,datacenter}")
	assert.Contains(t, collector.DiskReadIops.String(), "variableLabels: {id,name,datacenter,disk}")
	assert.Contains(t, collector.NetworkInPps.String(), "variableLabels: {id,name,datacenter,interface}")
}
