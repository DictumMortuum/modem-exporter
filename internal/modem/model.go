package modem

import "fmt"

type Stats struct {
	Uptime      int64 `json:"uptime"`
	Status     bool `json:"status"`
	CurrentUp   int   `json:"current_up"`
	CurrentDown int   `json:"current_down"`
	CRCUp       int   `json:"crc_up"`
	CRCDown     int   `json:"crc_down"`
}

// ToString method returns a string of the current statistics struct.
func (s *Stats) ToString() string {
	return fmt.Sprintf("%d / %d ", s.CurrentUp, s.CurrentDown)
}
