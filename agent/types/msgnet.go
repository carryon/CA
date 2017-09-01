package types

type MsgnetConfig struct {
	Logger struct {
		Level     string `json:"level"`
		Formatter string `json:"formatter"`
		Out       string `json:"out"`
	} `json:"logger"`
	Profiler struct {
		Port int `json:"port"`
	} `json:"profiler"`
	Router struct {
		ID                int         `json:"id"`
		Address           string      `json:"address"`
		AddressAutoDetect bool        `json:"addressAutoDetect"`
		Discovery         interface{} `json:"discovery"`
		Timeout           struct {
			Keepalive string `json:"keepalive"`
			Routers   string `json:"routers"`
			Network   struct {
				Routers string `json:"routers"`
				Peers   string `json:"peers"`
			} `json:"network"`
		} `json:"timeout"`
		Reconnect struct {
			Interval string `json:"interval"`
			Max      int    `json:"max"`
		} `json:"reconnect"`
	} `json:"router"`
	Report struct {
		On       bool   `json:"on"`
		ServerIP string `json:"serverIP"`
		Interval string `json:"interval"`
	} `json:"report"`
	NodeID     string `json:"node_id"`
	UpdateTime int    `json:"update_time"`
}
