package bgs

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

type SovereigntyMetrics struct {
	EventsProcessed    prometheus.CounterVec
	FilterLatency      prometheus.HistogramVec
	ActiveConnections  prometheus.GaugeVec
	CanadianEventsSent prometheus.Counter
	FilteredEvents     prometheus.Counter
}

var (
	sovereigntyMetricsInstance *SovereigntyMetrics
	sovereigntyMetricsOnce     sync.Once
)

// ResetSovereigntyMetricsForTesting resets the singleton for testing purposes
func ResetSovereigntyMetricsForTesting() {
	sovereigntyMetricsInstance = nil
	sovereigntyMetricsOnce = sync.Once{}
}

func NewSovereigntyMetrics() *SovereigntyMetrics {
	sovereigntyMetricsOnce.Do(func() {
		sovereigntyMetricsInstance = &SovereigntyMetrics{
			EventsProcessed: *promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "sovereignty_events_processed_total",
					Help: "Total number of events processed by sovereignty filter",
				},
				[]string{"result", "country"},
			),
			FilterLatency: *promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name: "sovereignty_filter_latency_seconds",
					Help: "Latency of geographic filtering in seconds",
				},
				[]string{"filter_type"},
			),
			ActiveConnections: *promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "sovereignty_active_connections",
					Help: "Number of active sovereign WebSocket connections",
				},
				[]string{"endpoint"},
			),
			CanadianEventsSent: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "sovereignty_canadian_events_sent_total",
					Help: "Total number of Canadian events sent via sovereign firehose",
				},
			),
			FilteredEvents: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "sovereignty_filtered_events_total",
					Help: "Total number of events filtered out by geographic filter",
				},
			),
		}
	})
	return sovereigntyMetricsInstance
}
