package modem

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DictumMortuum/modem-exporter/config"
	"github.com/DictumMortuum/modem-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

// Client struct is a PI-Hole client to request an instance of a PI-Hole ad blocker.
type Client struct {
	httpClient http.Client
	interval   time.Duration
	config     *config.Config
}

// NewClient method initializes a new PI-Hole client.
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

// Metrics scrapes pihole and sets them
func (c *Client) Metrics() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		stats, err := c.getStatistics()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		c.setMetrics(stats)

		log.Printf("New tick of statistics: %s", stats.ToString())
		promhttp.Handler().ServeHTTP(writer, request)
	}
}

func (c *Client) setMetrics(stats *Stats) {
	metrics.Uptime.WithLabelValues(c.config.Hostname).Set(float64(stats.Uptime))
	metrics.CurrentUp.WithLabelValues(c.config.Hostname).Set(float64(stats.CurrentUp))
	metrics.CurrentDown.WithLabelValues(c.config.Hostname).Set(float64(stats.CurrentDown))
	metrics.CRCUp.WithLabelValues(c.config.Hostname).Set(float64(stats.CRCUp))
	metrics.CRCDown.WithLabelValues(c.config.Hostname).Set(float64(stats.CRCDown))

	var isEnabled int = 0
	if stats.Status == true {
		isEnabled = 1
	}

	metrics.Status.WithLabelValues(c.config.Hostname).Set(float64(isEnabled))
}

func (c *Client) getStatistics() (*Stats, error) {
	retval := new(Stats)
	ip := c.config.Hostname
	ppp := "ip"
	ppp = strings.ToUpper(ppp)

	req, err := http.NewRequest("GET", "http://"+ip+"/comm/wan_cfg.sjs", nil)
	if err != nil {
		return nil, err
	}

	res1, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res1.Body.Close()
	if res1.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res1.StatusCode, res1.Status))
	}

	bodyBytes, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		return nil, err
	}

	vm := otto.New()

	_, err = vm.Run(string(bodyBytes))
	if err != nil {
		return nil, err
	}

	// data, err := vm.Run(ppp + "_ConnectionTable[0].TxBytes")
	// if err != nil {
	// 	return nil, err
	// }

	// retval.DataUp, _ = data.ToInteger()

	// data, err = vm.Run(ppp + "_ConnectionTable[0].RxBytes")
	// if err != nil {
	// 	return nil, err
	// }

	// retval.DataDown, _ = data.ToInteger()

	data, err := vm.Run(ppp + "_ConnectionTable[0].UpTime")
	if err != nil {
		return nil, err
	}
	uptime, _ := data.ToInteger()
	retval.Uptime = uptime

	data, err = vm.Run("GetWanDSLStatus()")
	if err != nil {
		return nil, err
	}

	// If the interface is bridged, it's not going to reset the uptime timer. So will do it manually here.
	dsl, _ := data.ToInteger()
	retval.Status = dsl == 1

	// t := time.Now().Add(time.Duration(-uptime) * time.Second)
	// retval.Date = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())

	res2, err := http.Get("http://" + ip + "/broadband/bd_dsl_detail.shtml?be=0&l0=2&l1=0&dtl=dt")
	if err != nil {
		return nil, err
	}
	defer res2.Body.Close()
	if res2.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res2.StatusCode, res2.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res2.Body)
	if err != nil {
		return nil, err
	}

	// doc.Find("td[key=PAGE_BD_DSL_DETAIL_MAXBDWIDTH] + td").Each(func(i int, s *goquery.Selection) {
	// 	current := strings.Split(s.Text(), "/")
	// 	retval.MaxUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
	// 	retval.MaxDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	// })

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_BDWIDTH] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.CurrentUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.CurrentDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
		// retval.InitialUp = retval.CurrentUp
		// retval.InitialDown = retval.CurrentDown
	})

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_CE] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.CRCUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.CRCDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	// doc.Find("td[key=PAGE_BD_DSL_DETAIL_FE] + td").Each(func(i int, s *goquery.Selection) {
	// 	current := strings.Split(s.Text(), "/")
	// 	retval.FECUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
	// 	retval.FECDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	// })

	firstScript := doc.Find("script[language=javascript]").First()

	vm = otto.New()

	_, err = vm.Run("function GetWanDSLStatus(){}")
	if err != nil {
		return nil, err
	}

	_, err = vm.Run(firstScript.Text())
	if err != nil {
		return nil, err
	}

	// snr, err := vm.Get("usNoiseMargin")
	// if err != nil {
	// 	return nil, err
	// }

	// retval.SNRUp, _ = snr.ToInteger()

	// snr, err = vm.Get("dsNoiseMargin")
	// if err != nil {
	// 	return nil, err
	// }

	// retval.SNRDown, _ = snr.ToInteger()

	return retval, nil
}
