package modem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ModemSpeedportPlus struct{}

func (m *ModemSpeedportPlus) GetStatistics(c *Client) (*Stats, error) {
	stats := new(Stats)

	req, err := http.NewRequest("GET", "http://"+c.config.Hostname+"/data/Status.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	type objtype struct {
		Type  string `json:"vartype"`
		Id    string `json:"varid"`
		Value string `json:"varvalue"`
	}

	rs := []objtype{}

	err = json.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	modem_config := map[string]string{}
	for _, item := range rs {
		modem_config[item.Id] = item.Value
	}

	// TODO: need to parse 95 days, 7 hours, 48 minutes, 55 seconds to unix timestamp
	stats.Uptime = 0
	stats.Status = modem_config["dsl_status"] == "online"
	stats.CurrentUp, _ = strconv.Atoi(modem_config["dsl_downstream"])
	stats.CurrentDown, _ = strconv.Atoi(modem_config["dsl_upstream"])
	stats.MaxUp, _ = strconv.Atoi(modem_config["dsl_max_downstream"])
	stats.MaxDown, _ = strconv.Atoi(modem_config["dsl_max_upstream"])
	stats.DataUp = 0
	stats.DataDown = 0
	// TODO: speedport only reports a single number of FEC errors, not two for up/down
	stats.FECUp, _ = strconv.Atoi(modem_config["dsl_fec_errors"])
	stats.FECDown = 0
	// TODO: speedport only reports a single number of CRC errors, not two for up/down
	stats.CRCUp, _ = strconv.Atoi(modem_config["dsl_crc_errors"])
	stats.CRCDown = 0
	// TODO: parse SNR correctly. Could snr be line attenuation here?
	snr := strings.Split(modem_config["dsl_snr"], "/")
	snr_up, _ := strconv.Atoi(strings.TrimSpace(snr[0]))
	snr_down, _ := strconv.Atoi(strings.TrimSpace(snr[1]))
	// fmt.Println(snr, snr_up, int64(snr_up))
	stats.SNRUp = int64(snr_up)
	stats.SNRDown = int64(snr_down)

	return stats, nil
}
