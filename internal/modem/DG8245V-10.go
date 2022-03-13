package modem

import (
	"github.com/ziutek/telnet"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DG8245V struct{}

const timeout = 10 * time.Second

var (
	re_max      = regexp.MustCompile(`Max:\s+Upstream rate = (\d+) Kbps, Downstream rate = (\d+) Kbps`)
	re_cur      = regexp.MustCompile(`Path:\s+\d+, Upstream rate = (\d+) Kbps, Downstream rate = (\d+) Kbps`)
	re_fec_down = regexp.MustCompile(`\nFECErrors:\s+(\d+)`)
	re_fec_up   = regexp.MustCompile(`ATUCFECErrors:\s+(\d+)`)
	re_crc_down = regexp.MustCompile(`\nCRCErrors:\s+(\d+)`)
	re_crc_up   = regexp.MustCompile(`ATUCCRCErrors:\s+(\d+)`)
	re_bytes    = regexp.MustCompile(`bytessent\s+= (\d+)\s+,bytesreceived\s+= (\d+)`)
	re_snr      = regexp.MustCompile(`display dsl snr up=([\d\.]+) down=([\d\.]+) success`)
	re_voip     = regexp.MustCompile(`Status\s+:Enable`)
)

func expect(t *telnet.Conn, d ...string) error {
	err := t.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

	err = t.SkipUntil(d...)
	if err != nil {
		return err
	}

	return nil
}

func sendln(t *telnet.Conn, s string) error {
	err := t.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'

	_, err = t.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func atoi64(s string) int64 {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return i
}

func atof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (m *DG8245V) GetStatistics(c *Client) (*Stats, error) {
	stats := new(Stats)

	t, err := telnet.Dial("tcp", c.config.Host+":23")
	if err != nil {
		return nil, err
	}
	defer t.Close()

	t.SetUnixWriteMode(true)
	var data []byte

	err = expect(t, "Login:")
	if err != nil {
		return nil, err
	}

	err = sendln(t, c.config.User)
	if err != nil {
		return nil, err
	}

	err = expect(t, "Password:")
	if err != nil {
		return nil, err
	}

	err = sendln(t, c.config.Pass)
	if err != nil {
		return nil, err
	}

	err = expect(t, "WAP>")
	if err != nil {
		return nil, err
	}

	err = sendln(t, "display xdsl connection status")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw := string(data)

	// TODO: need to parse On Line: 0 Days 3 Hour 17 Min 24 Sec to unix timestamp
	stats.Uptime = 0
	stats.Status = strings.Contains(raw, "Status: Up")

	refs := re_max.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.MaxUp = atoi(match[1])
		stats.MaxDown = atoi(match[2])
	}

	refs = re_cur.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.CurrentUp = atoi(match[1])
		stats.CurrentDown = atoi(match[2])
	}

	refs = re_crc_down.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.CRCDown = atoi(match[1])
	}

	refs = re_crc_up.FindAllStringSubmatch(raw, 1)
	if len(refs) > 0 {
		match := refs[0]
		stats.CRCUp = atoi(match[1])
	}

	refs = re_fec_down.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.FECDown = atoi(match[1])
	}

	refs = re_fec_up.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.FECUp = atoi(match[1])
	}

	err = sendln(t, "display xdsl statistics")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs = re_bytes.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.DataUp = atoi64(match[1])
		stats.DataDown = atoi64(match[2])
	}

	err = sendln(t, "display dsl snr")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs = re_snr.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.SNRUp = atof(match[1])
		stats.SNRDown = atof(match[2])
	}

	err = sendln(t, "display waninfo interface "+c.config.Voip)
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs = re_voip.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		stats.VoipStatus = true
	}

	return stats, nil
}
