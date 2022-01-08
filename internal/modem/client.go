package modem

import (
	"github.com/DictumMortuum/modem-exporter/config"
	"github.com/DictumMortuum/modem-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Client struct {
	httpClient http.Client
	config     *config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
		httpClient: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *Client) Metrics() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		stats, err := c.getStatistics()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		c.setMetrics(stats)

		promhttp.Handler().ServeHTTP(writer, request)
	}
}

func (c *Client) setMetrics(stats *Stats) {
	metrics.Uptime.WithLabelValues(c.config.Host).Set(float64(stats.Uptime))
	metrics.CurrentUp.WithLabelValues(c.config.Host).Set(float64(stats.CurrentUp))
	metrics.CurrentDown.WithLabelValues(c.config.Host).Set(float64(stats.CurrentDown))
	metrics.CRCUp.WithLabelValues(c.config.Host).Set(float64(stats.CRCUp))
	metrics.CRCDown.WithLabelValues(c.config.Host).Set(float64(stats.CRCDown))
	metrics.MaxUp.WithLabelValues(c.config.Host).Set(float64(stats.MaxUp))
	metrics.MaxDown.WithLabelValues(c.config.Host).Set(float64(stats.MaxDown))
	metrics.DataUp.WithLabelValues(c.config.Host).Set(float64(stats.DataUp))
	metrics.DataDown.WithLabelValues(c.config.Host).Set(float64(stats.DataDown))
	metrics.FECUp.WithLabelValues(c.config.Host).Set(float64(stats.FECUp))
	metrics.FECDown.WithLabelValues(c.config.Host).Set(float64(stats.FECDown))
	metrics.SNRUp.WithLabelValues(c.config.Host).Set(float64(stats.SNRUp))
	metrics.SNRDown.WithLabelValues(c.config.Host).Set(float64(stats.SNRDown))

	var isEnabled int = 0
	if stats.Status == true {
		isEnabled = 1
	}

	metrics.Status.WithLabelValues(c.config.Host).Set(float64(isEnabled))
}

func (c *Client) getStatistics() (*Stats, error) {
	switch c.config.Modem {
	case "TD5130":
		rs := ModemTD5130{}
		return rs.GetStatistics(c)
	case "SpeedportPlus":
		rs := ModemSpeedportPlus{}
		return rs.GetStatistics(c)
	default:
		return nil, nil
	}
}
