package modem

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ModemTD5130 struct{}

func (m *ModemTD5130) GetStatistics(c *Client) (*Stats, error) {
	retval := new(Stats)
	ip := c.config.Hostname
	ppp := "ip"
	ppp = strings.ToUpper(ppp)

	req, err := http.NewRequest("GET", "http://"+ip+"/comm/wan_cfg.sjs", nil)
	if err != nil {
		return nil, err
	}

	res1, err := c.httpClient.Do(req)
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

	data, err := vm.Run(ppp + "_ConnectionTable[0].TxBytes")
	if err != nil {
		return nil, err
	}

	retval.DataUp, _ = data.ToInteger()

	data, err = vm.Run(ppp + "_ConnectionTable[0].RxBytes")
	if err != nil {
		return nil, err
	}

	retval.DataDown, _ = data.ToInteger()

	data, err = vm.Run(ppp + "_ConnectionTable[0].UpTime")
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

	req, err = http.NewRequest("GET", "http://"+ip+"/broadband/bd_dsl_detail.shtml?be=0&l0=2&l1=0&dtl=dt", nil)
	if err != nil {
		return nil, err
	}

	res2, err := c.httpClient.Do(req)
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

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_MAXBDWIDTH] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.MaxUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.MaxDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_BDWIDTH] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.CurrentUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.CurrentDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_CE] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.CRCUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.CRCDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_FE] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.FECUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.FECDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

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

	snr, err := vm.Get("usNoiseMargin")
	if err != nil {
		return nil, err
	}

	retval.SNRUp, _ = snr.ToInteger()

	snr, err = vm.Get("dsNoiseMargin")
	if err != nil {
		return nil, err
	}

	retval.SNRDown, _ = snr.ToInteger()

	return retval, nil
}
