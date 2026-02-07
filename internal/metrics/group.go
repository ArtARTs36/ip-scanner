package metrics

import goMetrics "github.com/artarts36/go-metrics"

type Metrics struct {
	MMDB           *MMDB
	MMDBAutoUpdate *MMDBAutoUpdate
}

func NewMetrics(registry goMetrics.Registry) *Metrics {
	return &Metrics{
		MMDB:           NewMMDB(registry),
		MMDBAutoUpdate: NewMMDBAutoUpdate(registry),
	}
}
