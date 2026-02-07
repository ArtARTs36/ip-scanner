package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	goMetrics "github.com/artarts36/go-metrics"
)

type MMDBAutoUpdate struct {
	cycles         prometheus.Counter
	startedLoads   prometheus.Counter
	completedLoads *prometheus.CounterVec
	lastCycle      prometheus.Gauge
}

func NewMMDBAutoUpdate(registry goMetrics.Registry) *MMDBAutoUpdate {
	return &MMDBAutoUpdate{
		cycles: registry.NewCounter(prometheus.CounterOpts{
			Name: "mmdb_auto_update_cycles_total",
			Help: "MMDB Auto update: count of cycles",
		}),
		startedLoads: registry.NewCounter(prometheus.CounterOpts{
			Name: "mmdb_auto_update_started_loads_total",
			Help: "MMDB Auto update: count of started new db loads",
		}),
		completedLoads: registry.NewCounterVec(prometheus.CounterOpts{
			Name: "mmdb_auto_update_completed_loads_total",
			Help: "MMDB Auto update: count of completed new db loads",
		}, []string{"status"}),
		lastCycle: registry.NewGauge(prometheus.GaugeOpts{
			Name: "mmdb_auto_update_last_cycle",
			Help: "MMDB Auto update: time of last cycle",
		}),
	}
}

func (m *MMDBAutoUpdate) Describe(ch chan<- *prometheus.Desc) {
	m.cycles.Describe(ch)
	m.startedLoads.Describe(ch)
	m.completedLoads.Describe(ch)
	m.lastCycle.Describe(ch)
}

func (m *MMDBAutoUpdate) Collect(ch chan<- prometheus.Metric) {
	m.cycles.Collect(ch)
	m.startedLoads.Collect(ch)
	m.completedLoads.Collect(ch)
	m.lastCycle.Collect(ch)
}

func (m *MMDBAutoUpdate) IncCycles() {
	m.cycles.Inc()
}

func (m *MMDBAutoUpdate) IncStartedLoads() {
	m.startedLoads.Inc()
}

func (m *MMDBAutoUpdate) IncCompletedLoads(status bool) {
	st := "FAIL"
	if status {
		st = "OK"
	}

	m.completedLoads.WithLabelValues(st).Inc()
}

func (m *MMDBAutoUpdate) UpdateLastCycle() {
	m.lastCycle.SetToCurrentTime()
}
