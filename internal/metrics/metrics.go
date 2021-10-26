package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	Uptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "uptime",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	CurrentUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "current_up",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	CurrentDown = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "current_down",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	CRCUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "crc_up",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	CRCDown = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "crc_down",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	Status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "enabled",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)
)

// Init initializes all Prometheus metrics made available by PI-Hole exporter.
func Init() {
	initMetric("modem_uptime", Uptime)
	initMetric("modem_current_up", CurrentUp)
	initMetric("modem_current_down", CurrentDown)
	initMetric("modem_crc_up", CRCUp)
	initMetric("modem_crc_down", CRCDown)
	initMetric("enabled", Status)
}

func initMetric(name string, metric *prometheus.GaugeVec) {
	prometheus.MustRegister(metric)
	log.Printf("New Prometheus metric registered: %s", name)
}
