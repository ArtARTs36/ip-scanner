package metrics

import (
	"fmt"
	"strings"
	"time"

	goMetrics "github.com/artarts36/go-metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type MMDB struct {
	actualTime prometheus.Gauge
	info       *prometheus.GaugeVec
	size       prometheus.Gauge
}

func NewMMDB(registry goMetrics.Registry) *MMDB {
	return &MMDB{
		actualTime: registry.NewGauge(prometheus.GaugeOpts{
			Name: "mmdb_actual_time",
			Help: "Actual time (epoch)",
		}),
		info: registry.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mmdb_info",
			Help: "Information about MMDB",
		}, []string{"type", "build_at", "languages", "size"}),
		size: registry.NewGauge(prometheus.GaugeOpts{
			Name: "mmdb_size",
			Help: "Size of .mmdb",
		}),
	}
}

func (m *MMDB) Scrap(actualTime time.Time, dbType string, languages []string, size int64) {
	m.actualTime.Set(float64(actualTime.Unix() / 1e9)) //nolint:mnd //not need

	m.info.WithLabelValues(
		dbType,
		actualTime.String(),
		strings.Join(languages, ","),
		fmt.Sprintf("%d", size),
	).Set(1)

	m.size.Set(float64(size))
}

func (m *MMDB) Describe(ch chan<- *prometheus.Desc) {
	m.actualTime.Describe(ch)
	m.info.Describe(ch)
}

func (m *MMDB) Collect(ch chan<- prometheus.Metric) {
	m.actualTime.Collect(ch)
	m.info.Collect(ch)
}
