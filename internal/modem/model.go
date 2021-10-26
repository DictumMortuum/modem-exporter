package modem

type Stats struct {
	Uptime      int64 `json:"uptime"`
	Status      bool  `json:"status"`
	CurrentUp   int   `json:"current_up"`
	CurrentDown int   `json:"current_down"`
	MaxUp       int   `json:"max_up"`
	MaxDown     int   `json:"max_down"`
	DataUp      int64 `json:"data_up"`
	DataDown    int64 `json:"data_down"`
	FECUp       int   `json:"fec_up"`
	FECDown     int   `json:"fec_down"`
	CRCUp       int   `json:"crc_up"`
	CRCDown     int   `json:"crc_down"`
	SNRUp       int64 `json:"snr_up"`
	SNRDown     int64 `json:"snr_down"`
}

type Statable interface {
	GetStatistics(*Client) (*Stats, error)
}
