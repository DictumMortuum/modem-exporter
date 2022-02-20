package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

var (
	Uptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "uptime",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	Status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "status",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	VoipStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "voip_status",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	MaxUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "max_up",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	MaxDown = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "max_down",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	DataUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "data_up",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	DataDown = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "data_down",
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

	FECUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "fec_up",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	FECDown = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "fec_down",
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

	SNRUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "snr_up",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)

	SNRDown = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "snr_down",
			Namespace: "modem",
		},
		[]string{"hostname"},
	)
)

func Init() {
	initMetric("modem_uptime", Uptime)
	initMetric("modem_status", Status)
	initMetric("modem_voip_status", VoipStatus)
	initMetric("modem_current_up", CurrentUp)
	initMetric("modem_current_down", CurrentDown)
	initMetric("modem_max_up", MaxUp)
	initMetric("modem_max_down", MaxDown)
	initMetric("modem_data_up", DataUp)
	initMetric("modem_data_down", DataDown)
	initMetric("modem_fec_up", FECUp)
	initMetric("modem_fec_down", FECDown)
	initMetric("modem_crc_up", CRCUp)
	initMetric("modem_crc_down", CRCDown)
	initMetric("modem_snr_up", SNRUp)
	initMetric("modem_snr_down", SNRDown)
}

func initMetric(name string, metric *prometheus.GaugeVec) {
	prometheus.MustRegister(metric)
	log.Printf("New Prometheus metric registered: %s", name)
}
